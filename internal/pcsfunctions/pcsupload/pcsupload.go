// Package pcsupload 上传包
package pcsupload

import (
	"github.com/wuchongde/BaiduPCS-Go/pcsverbose"
)

const (
	UploadingFileName = "pcs_uploading.json"
)

var (
	pcsUploadVerbose = pcsverbose.New("PCSUPLOAD")
)
