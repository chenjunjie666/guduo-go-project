package release_time

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestDoubanReleaseTime(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Baidu.GetBaikeUrl()[0]
	baiduReleaseTime(url.Url, url.ShowId)
}
