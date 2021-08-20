package fans

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestWeiboFans(t *testing.T) {
	core.Init()

	wg.Add(1)
	row := storage.Weibo.GetDetailUrl()[0]
	weiboFans(row.Url, row.ShowId)
	// weiboFans("https://weibo.com/huaweiweibo", 100)
}
