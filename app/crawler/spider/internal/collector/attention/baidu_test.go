package attention

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestBaiduAttention(t *testing.T) {
	core.Init()
	wg.Add(1)
	url := storage.Baidu.GetDetailUrl()[0]
	baiduAttention(url.Url, url.ShowId)
}