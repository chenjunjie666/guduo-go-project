package article_num

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestWeiboArticleNum(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Weibo.GetFetchArticleUrl()[0]
	weiboArticleNum(url.Url, url.ID)
}
