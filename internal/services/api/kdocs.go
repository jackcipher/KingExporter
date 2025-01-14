package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"KingExporter/internal/global"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
)

const KDocsSID = "wps_sid"

type KDocsApi struct {
	baseHost  string
	driveHost string
	sid       string
	client    *resty.Client
}

func NewKDocsApi(baseHost string, driveHost string, sid string) *KDocsApi {
	c := new(KDocsApi)
	c.baseHost = baseHost
	c.driveHost = driveHost
	c.sid = sid
	c.client = resty.New()

	return c
}

func (c *KDocsApi) Req(debug ...bool) *resty.Request {
	isDebug := false
	if len(debug) > 0 && debug[0] {
		isDebug = true
	}
	return c.client.SetDebug(isDebug).R().SetCookie(&http.Cookie{
		Name:  KDocsSID,
		Value: c.sid,
	}).SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36").
		SetHeader("Referer", c.baseHost).SetHeader("Origin", c.baseHost)
}

type UserInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Status string `json:"status"`
}

func (c *KDocsApi) UserInfo() (*UserInfo, error) {
	endpoint := fmt.Sprintf("%s/api/v3/userinfo", c.driveHost)
	resp, err := c.Req().Get(endpoint)
	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] get user info error: %v", err))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] get user info empty response")
		return nil, err
	}

	var userInfo UserInfo
	if err := json.Unmarshal(resp.Body(), &userInfo); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal user info error: %v", err))
		return nil, err
	}

	return &userInfo, nil
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c *KDocsApi) GetGroups() ([]Group, error) {
	var respData struct {
		Groups []Group `json:"groups"`
	}
	endpoint := fmt.Sprintf("%s/api/v3/groups", c.driveHost)
	resp, err := c.Req().Get(endpoint)
	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] GetGroups failed: %s", err.Error()))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] GetGroups empty response")
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &respData); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal groups error: %s", err))
		return nil, err
	}
	return respData.Groups, nil
}

type File struct {
	ID       int    `json:"id"`
	ParentID int    `json:"parentid"`
	FName    string `json:"fname"`
	FSize    int    `json:"fsize"`
	FType    string `json:"ftype"`
}

func (c *KDocsApi) Files(groupID, parentID int) ([]File, error) {
	var data struct {
		Files []File `json:"files"`
	}
	endpoint := fmt.Sprintf("%s/api/v5/groups/%d/files?parentid=%d&offset=0&count=20000", c.driveHost, groupID, parentID)
	resp, err := c.Req().Get(endpoint)

	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] Files failed: %s", err.Error()))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] Files empty response")
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal Files error: %s", err))
		return nil, err
	}
	return data.Files, nil
}

type DownloadItem struct {
	DownloadUrl string `json:"download_url"`
	Url         string `json:"url"`
	// 金山把字段打错了 应该是 fsize
	Size int `json:"fize"`
}

type PDFDownloadItem struct {
	Size int    `json:"fsize"`
	Url  string `json:"url"`
}

func (c *KDocsApi) GetPDFDownloadUrl(groupID, fileID int) (*PDFDownloadItem, error) {
	var data PDFDownloadItem
	endpoint := fmt.Sprintf("%s/api/v5/groups/%d/files/%d/download?support_checksums=md5", c.driveHost, groupID, fileID)

	resp, err := c.Req().Get(endpoint)
	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] GetPDFDownloadUrl failed: %s", err.Error()))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] GetPDFDownloadUrl empty response")
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal GetPDFDownloadUrl error: %s", err))
		return nil, err
	}

	return &data, nil
}

func (c *KDocsApi) GetDownloadUrl(fileID int) (*DownloadItem, error) {
	var data DownloadItem
	endpoint := fmt.Sprintf("%s/api/v3/office/file/%d/download", c.baseHost, fileID)

	resp, err := c.Req().Get(endpoint)
	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] GetDownloadUrl failed: %s", err.Error()))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] GetDownloadUrl empty response")
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal GetDownloadUrl error: %s", err))
		return nil, err
	}

	return &data, nil
}

type ExportPreload struct {
	TaskID   string `json:"task_id"`
	TaskType string `json:"task_type"`
}

type ExportResult struct {
	Status string `json:"status"`
	Data   struct {
		Key string `json:"key"`
		Url string `json:"url"`
	} `json:"data"`
}

func (c *KDocsApi) PreloadExport(fileID int, filename string) (*ExportPreload, error) {
	var data ExportPreload
	format := c.GetFormat(filename)
	endpoint := fmt.Sprintf("%s/api/v3/office/file/%d/export/%s/preload", c.baseHost, fileID, format)
	resp, err := c.Req().SetBody(map[string]string{
		"ver":    "3",
		"format": format,
	}).Post(endpoint)

	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] PreloadExport failed: %s", err.Error()))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] PreloadExport empty response")
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal PreloadExport error: %s", err))
		return nil, err
	}

	return &data, nil
}

func (c *KDocsApi) ExportProgress(fileID int, taskID, taskType, format string) (*ExportResult, error) {
	var data ExportResult
	endpoint := fmt.Sprintf("%s/api/v3/office/file/%d/export/%s/result", c.baseHost, fileID, format)
	resp, err := c.Req().SetBody(map[string]string{
		"task_id":   taskID,
		"task_type": taskType,
	}).Post(endpoint)
	if err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] ExportProgress failed: %s", err.Error()))
		return nil, err
	}

	if resp == nil {
		global.Log.Error("[KDocsApi] ExportProgress empty response")
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		global.Log.Error(fmt.Sprintf("[KDocsApi] unmarshal ExportProgress error: %s", err))
		return nil, err
	}

	return &data, nil
}

func (c *KDocsApi) GetFormat(fName string) string {
	ext := path.Ext(fName)
	return lo.If(ext == ".otl", "docx").Else("xlsx")
}
