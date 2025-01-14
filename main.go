package main

import (
	"flag"
	"os"

	"KingExporter/internal/services/kdocs"
	"KingExporter/pkg/display"
)

type flags struct {
	silent      bool
	downloadDir string
	exportAll   bool
	groupID     int
	sid         string
}

func parseFlags() *flags {
	f := &flags{}

	flag.BoolVar(&f.silent, "s", false, "开启静默模式")
	flag.StringVar(&f.downloadDir, "download_dir", "", "下载文件的目录")
	flag.BoolVar(&f.exportAll, "A", false, "是否导出所有的文档，包括个人文档及团队文档")
	flag.IntVar(&f.groupID, "group_id", 0, "导出指定空间的文档")
	flag.StringVar(&f.sid, "sid", "", "金山文档的会话 ID")

	flag.Parse()
	return f
}

func main() {
	f := parseFlags()
	e := kdocs.NewExporter(f.sid, kdocs.ExportOptions{
		DownloadDir: f.downloadDir,
		SilentMode:  f.silent,
		ExportAll:   f.exportAll,
		GroupID:     f.groupID,
	})

	dir := e.Export()
	if !f.silent {
		waitForKeyPress(dir)
	}
}

func waitForKeyPress(dir string) {
	display.Print("下载地址： \n\t%s\n\n程序处理结束，请按任意键退出...", dir)
	// 为了防止 Windows 机器 直接点开 exe  程序自动退出  看不到下载目录
	os.Stdin.Read(make([]byte, 1))
}
