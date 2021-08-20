package danmaku

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestBilibiliDanmakuContent(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Bilibili.GetDetailUrl()[0]
	bilibiliDanmakuContent(url.Url, url.ShowId)
}
