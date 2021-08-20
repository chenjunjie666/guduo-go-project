package upload

import (
	"bytes"
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"guduo/pkg/errors"
	"strings"
	"time"
)

var expired uint64
var upToken string
var uriPath = "test_upload/"
var bucket = "test"
var domain = "http://127.0.0.1/"

func Image(data []byte, filename string, nameSalt ...string) (string, error) {
	if int(expired)-int(time.Now().Unix()) < 60 {
		genToken()
	}

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	dataLen := int64(len(data))

	if len(nameSalt) > 0 {
		filename = strings.Join(nameSalt, "_") + "_" + filename
	}

	path := uriPath + filename
	err := formUploader.Put(context.Background(), &ret, upToken, path, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return "", errors.CmsError("上传文件失败")
	}
	return domain + path, nil
}

func genToken() {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	expired = uint64(time.Now().Unix()) + 86400*365
	putPolicy.Expires = 86400 * 365 // 一年有效期
	ak := "zp_FGY16eBCRZTuKtmWjksMl_Bg0u3OvkLv2VGF9"
	sk := "ST7PT4N-Jqer8Q8ndbGqYirJd-_79jO3VWZpU4bh"
	mac := qbox.NewMac(ak, sk)
	upToken = putPolicy.UploadToken(mac)

}
