package article_num_actor

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestWechatArticleNum(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Wechat.GetArticleNumUrl()[0]
	wechatArticleNum(url.Url, url.ID)
}

