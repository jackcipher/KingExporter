package kdocs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"KingExporter/internal/global"
	"KingExporter/internal/services/api"
	"KingExporter/pkg/display"
	"KingExporter/pkg/utils"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const (
	ApiHostBase  = "https://www.kdocs.cn"
	ApiHostDrive = "https://drive.kdocs.cn"
)

const (
	NumWorkerDownload = 20
	NumWorkerPreload  = 10
	MaxRetries        = 3
)

type state struct {
	downloadCh  chan DownloadJob
	preloadCh   chan PreloadJob
	downloadWg  *sync.WaitGroup
	preloadWg   *sync.WaitGroup
	workerWg    *sync.WaitGroup
	downloadDir string
}

type Exporter struct {
	silent      bool
	downloadDir string
	sid         string
	api         *api.KDocsApi
	exportAll   bool
	groupID     int
}

type ExportOptions struct {
	DownloadDir string
	SilentMode  bool
	ExportAll   bool
	GroupID     int
}

func NewExporter(sid string, options ExportOptions) *Exporter {
	e := &Exporter{
		silent:      options.SilentMode,
		downloadDir: options.DownloadDir,
		groupID:     options.GroupID,
		exportAll:   options.ExportAll,
		sid:         sid,
	}

	e.Check()
	return e
}

func (e *Exporter) EnableSilent(silent bool) {
	e.silent = silent
}

func (e *Exporter) getTempDir() (string, error) {
	dir := path.Join(os.TempDir(), "kdocs-files")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		err := fmt.Errorf("创建临时文件失败: %w", err)
		global.Log.Error(err.Error())
		return "", err
	}

	return dir, nil
}

func (e *Exporter) exportGroup(groupID int, name ...string) {
	st := &state{
		downloadCh: make(chan DownloadJob, NumWorkerDownload),
		preloadCh:  make(chan PreloadJob, NumWorkerPreload),
		downloadWg: &sync.WaitGroup{},
		preloadWg:  &sync.WaitGroup{},
		workerWg:   &sync.WaitGroup{},
	}
	st.downloadDir = e.downloadDir
	if len(name) > 0 {
		st.downloadDir = path.Join(e.downloadDir, name[0])
	}

	for i := 0; i < NumWorkerPreload; i++ {
		st.workerWg.Add(1)
		go e.preloadWorker(i, st)
	}

	for i := 0; i < NumWorkerDownload; i++ {
		st.workerWg.Add(1)
		go e.downloadWorker(i, st)
	}

	// DFS 遍历目录
	if err := e.processFolder(groupID, 0, "", st); err != nil {
		panic(err)
	}

	// 等待所有转码任务结束
	st.preloadWg.Wait()
	// 所有转码任务加入 channel 后，关闭 channel
	close(st.preloadCh)
	// 等待所有转码任务结束
	st.downloadWg.Wait()
	// 所有的下载任务写入到 downloadCh 中，关闭 channel
	close(st.downloadCh)

	// 等待所有的 worker 结束
	st.workerWg.Wait()
}

func (e *Exporter) Export() string {
	groups, err := e.api.GetGroups()
	if err != nil {
		err = fmt.Errorf("获取我的云文件及团队 group 失败: %w", err)
		global.Log.Error(err.Error())
		display.Exit(1, err.Error())
	}

	dir := e.downloadDir
	if e.exportAll {
		wg := sync.WaitGroup{}
		for _, v := range groups {
			go func() {
				wg.Add(1)
				subDir := path.Join(e.downloadDir, v.Name)
				if err := os.MkdirAll(subDir, os.ModePerm); err != nil {
					err = fmt.Errorf("创建 %s group 失败: %w", v.Name, err)
					global.Log.Error(err.Error())
					return
				}
				e.exportGroup(v.ID, v.Name)
				fmt.Printf("✅ 团队 %s 文档导出完成\n", v.Name)
				wg.Done()
			}()
		}
		wg.Wait()
	} else if e.groupID > 0 {
		found := false
		for _, v := range groups {
			if v.ID == e.groupID {
				e.exportGroup(e.groupID, v.Name)
				found = true
				break
			}
		}
		if !found {
			display.Exit(1, "groupID: %d 不存在", e.groupID)
		}
	} else {
		for _, v := range groups {
			if v.Type == "special" || v.Type == "corpspecial" {
				e.exportGroup(v.ID, v.Name)
			}
		}
	}

	return dir
}

func (e *Exporter) processFile(f api.File, groupID int, relativePath string, st *state) error {
	fullPath := filepath.Join(st.downloadDir, relativePath, f.FName)
	dirPath := filepath.Dir(fullPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
	}

	ext := filepath.Ext(f.FName)
	if lo.Contains([]string{".docx", ".pptx", ".doc", ".ppt", ".xls", ".xlsx"}, ext) {
		item, err := e.api.GetDownloadUrl(f.ID)
		if err != nil {
			global.Log.Error(e.logError("获取文件见地址失败", err, f, groupID))
			return err
		}
		st.downloadWg.Add(1)
		st.downloadCh <- DownloadJob{Url: item.Url, FullPath: fullPath}
	} else if ext == ".pdf" {
		item, err := e.api.GetPDFDownloadUrl(groupID, f.ID)
		if err != nil {
			global.Log.Error(e.logError("获取 PDF 下载地址失败", err, f, groupID))
			return err
		}
		st.downloadWg.Add(1)
		st.downloadCh <- DownloadJob{Url: item.Url, FullPath: fullPath}
	} else if lo.Contains([]string{".otl", ".ksheet"}, ext) {
		fullPath = lo.If(ext == ".otl", utils.ReplaceExt(fullPath, ".docx")).Else(utils.ReplaceExt(fullPath, ".xlsx"))
		st.preloadWg.Add(1)
		st.preloadCh <- PreloadJob{
			File:       f,
			GroupID:    groupID,
			FullPath:   fullPath,
			RetryCount: 0,
			MaxRetries: MaxRetries,
		}
	}

	return nil
}

func (e *Exporter) processFolder(groupID int, folderID int, relativePath string, st *state) error {
	files, err := e.api.Files(groupID, folderID)
	if err != nil {
		return fmt.Errorf("获取目录文件失败 folderID %d: %v", folderID, err)
	}

	for _, file := range files {
		if file.FType == "folder" {
			newPath := filepath.Join(relativePath, file.FName)
			if err := e.processFolder(groupID, file.ID, newPath, st); err != nil {
				global.Log.Error(fmt.Sprintf("处理文件夹失败 %s: %v", file.FName, err))
				continue
			}
		} else {
			if err := e.processFile(file, groupID, relativePath, st); err != nil {
				global.Log.Error(fmt.Sprintf("处理文件失败 %s: %v", file.FName, err))
				continue
			}
		}
	}

	return nil
}

func (e *Exporter) downloadWorker(id int, st *state) {
	defer st.workerWg.Done()
	for {
		select {
		case job, ok := <-st.downloadCh:
			if !ok {
				return
			}

			resp, err := resty.New().R().Head(job.Url)
			if err != nil {
				err = fmt.Errorf("[Download #[%d] Failed to get download info: %w", id, err)
				global.Log.Error("查看下载信息失败", err.Error())
			}
			size := cast.ToInt64(resp.Header().Get("Content-Length"))
			slow := lo.If(size > 100<<20, "⚠️ slow").Else("")
			fmt.Printf("⏬ Downloading to %s fileSize: %s %s\n", job.FullPath, display.FormatBytes(size), slow)
			_, err = resty.New().R().SetOutput(job.FullPath).Get(job.Url)
			if err != nil {
				global.Log.Error(fmt.Sprintf("[Download #%d] Failed to download %s to %s: %s", id, job.Url, job.FullPath, err.Error()))
			}
			st.downloadWg.Done()
		}
	}
}
