package short_comment

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestDoubanShortComment(t *testing.T) {
	core.Init()
	wg.Add(1)

	url := storage.Douban.GetDetailUrl()[0]
	doubanShortComment(url.Url, url.ShowId)
}
