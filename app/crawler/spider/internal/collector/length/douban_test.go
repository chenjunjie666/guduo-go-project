package length

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestDoubanLength(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Douban.GetDetailUrl()[0]
	doubanLength(url.Url, url.ShowId)
}
