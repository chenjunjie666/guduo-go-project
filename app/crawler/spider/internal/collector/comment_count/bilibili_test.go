package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestBilibiliCommentCount(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Bilibili.GetDetailUrl()[0]
	bilibiliCommentCount(url.Url, url.ShowId)
	// bilibiliCommentCount("https://www.bilibili.com/bangumi/play/ep401327", 100) 国创
}
