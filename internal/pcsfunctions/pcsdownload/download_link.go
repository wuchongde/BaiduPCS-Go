package pcsdownload

import (
	"github.com/wuchongde/BaiduPCS-Go/baidupcs"
	"github.com/wuchongde/BaiduPCS-Go/internal/pcsconfig"
	"net/url"
)

func GetLocateDownloadLinks(pcs *baidupcs.BaiduPCS, pcspath string) (dlinks []*url.URL, err error) {
	dInfo, pcsError := pcs.LocateDownload(pcspath)
	if pcsError != nil {
		return nil, pcsError
	}

	us := dInfo.URLStrings(pcsconfig.Config.EnableHTTPS)
	if len(us) == 0 {
		return nil, ErrDlinkNotFound
	}

	return us, nil
}
