package storage_test

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestWeibo_GetFetchArticleUrl(t *testing.T) {
	core.Init()

	u := storage.Weibo.GetFetchArticleUrl()
	fmt.Println(u)
}
