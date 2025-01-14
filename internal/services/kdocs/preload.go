package kdocs

import (
	"fmt"
	"time"

	"KingExporter/internal/global"
)

func (e *Exporter) preloadWorker(id int, st *state) {
	defer st.workerWg.Done()
	for {
		select {
		case job, ok := <-st.preloadCh:
			if !ok {
				return
			}
			fmt.Printf("⌛️ Preload export %s\n", job.File.FName)
			err := e.handlePreload(job, st)
			if err != nil && job.RetryCount < job.MaxRetries {
				job.RetryCount++
				time.Sleep(time.Second)
				st.preloadCh <- job
			} else if err != nil {
				global.Log.Error(fmt.Sprintf("[Preload #%d] Failed to process %s after %d retries: %s",
					id, job.File.FName, job.RetryCount, err))
			}
			st.preloadWg.Done()
		}
	}
}

func (e *Exporter) handlePreload(job PreloadJob, st *state) error {
	data, err := e.api.PreloadExport(job.File.ID, job.File.FName)
	if err != nil {
		global.Log.Error(e.logError("预导出文件失败", err, job.File, job.GroupID))
		return err
	}
	if data.TaskID == "" {
		err = fmt.Errorf("预导出文件，taskID 为空")
		global.Log.Error("预导出文件失败", err)
		return err
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-timeout:
			global.Log.Error("转码导出超时", job.File.FName, job.File.FSize)
			return fmt.Errorf("转码导出失败")
		case <-ticker.C:
			result, err := e.api.ExportProgress(
				job.File.ID,
				data.TaskID,
				data.TaskType,
				e.api.GetFormat(job.File.FName),
			)
			if err != nil {
				global.Log.Error(fmt.Sprintf("获取导出进度失败: %s", err), job.File.FSize, job.File.FName)
				continue
			}
			if result.Status == "finished" {
				st.downloadWg.Add(1)
				st.downloadCh <- DownloadJob{
					Url:      result.Data.Url,
					FullPath: job.FullPath,
				}
				return nil
			}
		}
	}
}
