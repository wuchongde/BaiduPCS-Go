package baidupcs

import (
	"errors"
	"github.com/wuchongde/BaiduPCS-Go/baidupcs/pcserror"
	"strings"
)

type (
	// ShareOption 分享可选项
	ShareOption struct {
		Password string // 密码
		Period   int    // 有效期
	}

	// Shared 分享信息
	Shared struct {
		Link    string `json:"link"`
		ShareID int64  `json:"shareid"`
	}

	// ShareRecordInfo 分享信息
	ShareRecordInfo struct {
		ShareID         int64   `json:"shareId"`
		FsIds           []int64 `json:"fsIds"`
		Passwd          string  `json:"-"` // 这个字段已经没有用了, 需要从ShareSURLInfo中获取
		Shortlink       string  `json:"shortlink"`
		Status          int     `json:"status"`          // 状态
		Public          int     `json:"public"`          // 是否为公开分享
		TypicalCategory int     `json:"typicalCategory"` // 文件类型
		TypicalPath     string  `json:"typicalPath"`
	}

	shareSURLInfo struct {
		*pcserror.PanErrorInfo
		*ShareSURLInfo
	}

	// ShareSURLInfo 分享的子信息
	ShareSURLInfo struct {
		Pwd      string `json:"pwd"` // 新密码
		ShortURL string `json:"shorturl"`
	}

	// ShareRecordInfoList 分享信息列表
	ShareRecordInfoList []*ShareRecordInfo

	sharePSetJSON struct {
		*Shared
		*pcserror.PanErrorInfo
	}

	shareListJSON struct {
		List ShareRecordInfoList `json:"list"`
		*pcserror.PanErrorInfo
	}
)

var (
	// ErrShareLinkNotFound 未找到分享链接
	ErrShareLinkNotFound = errors.New("未找到分享链接")
)

// ShareSet 分享文件
func (pcs *BaiduPCS) ShareSet(paths []string, option *ShareOption) (s *Shared, pcsError pcserror.Error) {
	if option == nil {
		option = &ShareOption{}
	}

	dataReadCloser, pcsError := pcs.PrepareSharePSet(paths, option.Period)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationShareSet)
	jsonData := sharePSetJSON{
		Shared:       &Shared{},
		PanErrorInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationShareSet, dataReadCloser, &jsonData)
	if pcsError != nil {
		return
	}

	if jsonData.Link == "" {
		errInfo.ErrType = pcserror.ErrTypeOthers
		errInfo.Err = ErrShareLinkNotFound
		return nil, errInfo
	}

	return jsonData.Shared, nil
}

// ShareCancel 取消分享
func (pcs *BaiduPCS) ShareCancel(shareIDs []int64) (pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareShareCancel(shareIDs)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	pcsError = pcserror.DecodePanJSONError(OperationShareCancel, dataReadCloser)
	return
}

// ShareList 列出分享列表
func (pcs *BaiduPCS) ShareList(page int) (records ShareRecordInfoList, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareShareList(page)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationShareList)
	jsonData := shareListJSON{
		List:         records,
		PanErrorInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationShareList, dataReadCloser, &jsonData)
	if pcsError != nil {
		// json解析错误
		if pcsError.GetErrType() == pcserror.ErrTypeJSONParseError {
			// 服务器更改, List为空时变成{}, 导致解析错误
			if strings.Contains(pcsError.GetError().Error(), `"list":{}`) {
				// 返回空列表
				return jsonData.List, nil
			}
		}
		return
	}

	if jsonData.List == nil {
		errInfo.ErrType = pcserror.ErrTypeOthers
		errInfo.Err = errors.New("shared list is nil")
		return nil, errInfo
	}

	return jsonData.List, nil
}

//ShareSURLInfo 获取分享的详细信息, 包含密码
func (pcs *BaiduPCS) ShareSURLInfo(shareID int64) (info *ShareSURLInfo, pcsError pcserror.Error) {
	dataReadCloser, pcsError := pcs.PrepareShareSURLInfo(shareID)
	if pcsError != nil {
		return
	}

	defer dataReadCloser.Close()

	errInfo := pcserror.NewPanErrorInfo(OperationShareSURLInfo)

	jsonData := shareSURLInfo{
		PanErrorInfo: errInfo,
	}

	pcsError = pcserror.HandleJSONParse(OperationShareList, dataReadCloser, &jsonData)
	if pcsError != nil {
		// json解析错误
		return
	}

	// 去掉0
	if jsonData.Pwd == "0" {
		jsonData.Pwd = ""
	}

	return jsonData.ShareSURLInfo, nil
}
