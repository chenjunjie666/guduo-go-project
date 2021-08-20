package rating_num

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestDoubanRatingNum(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Douban.GetDetailUrl()[0]
	doubanRatingNum(url.Url, url.ShowId)
}
