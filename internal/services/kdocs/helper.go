package kdocs

import (
	"errors"
	"fmt"
	"os"

	"KingExporter/internal/global"
	"KingExporter/internal/services/api"
	"KingExporter/pkg/display"
)

// scanInput is a helper function to handle user input
func (e *Exporter) scanInput(target *string) {
	_, err := fmt.Scanln(target)
	if err != nil {
		global.Log.Error("获取用户输入失败", err)
	}
}

func (e *Exporter) logError(msg string, err error, f api.File, groupID int) string {
	return fmt.Sprintf("[%s]: %s GroupID: %d fileID: %d fileName: %s", msg, err, groupID, f.ID, f.FName)
}

// validateUserAccess verifies the user's credentials and access
func (e *Exporter) validateUserAccess() error {
	e.api = api.NewKDocsApi(ApiHostBase, ApiHostDrive, e.sid)

	for {
		userinfo, err := e.api.UserInfo()
		if err != nil {
			if e.silent {
				return fmt.Errorf("获取用户信息失败: %w", err)
			}
			display.PrintInput("获取用户信息失败，请输入正确的 sid")
			e.scanInput(&e.sid)
			e.api = api.NewKDocsApi(ApiHostBase, ApiHostDrive, e.sid)
			continue
		}

		if !e.silent {
			display.Print("当前登录用户: \n\t%s", userinfo.Name)
		}
		break
	}
	return nil
}

// validateSID ensures a valid session ID is provided
func (e *Exporter) validateSID() error {
	if e.sid != "" {
		return nil
	}

	if e.silent {
		return errors.New("静默模式下请通过 --sid 指定金山文档的会话 ID")
	}

	display.PrintInput("请输入金山文档的会话 ID (sid):")
	e.scanInput(&e.sid)
	return nil
}

// setupDownloadDirectory configures and validates the download directory
func (e *Exporter) setupDownloadDirectory() error {
	if e.downloadDir == "" {
		defaultDir, err := e.getTempDir()
		tip := "获取临时文件夹失败，请输入下载地址"
		if err != nil {
			global.Log.Warn("获取临时文件夹失败", err)
		} else {
			tip = fmt.Sprintf("请输入下载地址 (%s)", defaultDir)
		}
		if !e.silent {
			display.PrintInput(tip)
			e.scanInput(&e.downloadDir)
		}

		if e.downloadDir == "" {
			e.downloadDir = defaultDir
		}
	}

	return e.validateDirectory()
}

// validateDirectory ensures the download directory exists and is actually a directory
func (e *Exporter) validateDirectory() error {
	for {
		f, err := os.Stat(e.downloadDir)
		if err != nil {
			if e.silent {
				return fmt.Errorf("下载目录不合法: %w", err)
			}
			display.PrintInput("下载地址目录设置不正确: %s", err.Error())
			e.scanInput(&e.downloadDir)
			continue
		}

		if !f.IsDir() {
			if e.silent {
				return errors.New("指定的下载路径非文件夹，请检查后重试")
			}
			display.PrintInput("你选择的路径非文件夹，请输入正确的文件夹路径")
			e.scanInput(&e.downloadDir)
			continue
		}
		break
	}
	return nil
}

func (e *Exporter) Check() {
	if err := e.validateSID(); err != nil {
		err = fmt.Errorf("获取会话信息失败: %w", err)
		global.Log.Error(err.Error())
		display.Exit(1, err.Error())
	}

	if err := e.setupDownloadDirectory(); err != nil {
		err = fmt.Errorf("设置云文件下载目录失败: %w", err)
		global.Log.Error(err.Error())
		display.Exit(1, err.Error())
	}

	if err := e.validateUserAccess(); err != nil {
		err = fmt.Errorf("获取用户信息失败: %w", err)
		global.Log.Error(err.Error())
		display.Exit(1, err.Error())
	}
}
