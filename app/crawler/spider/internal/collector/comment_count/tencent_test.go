package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestTencentCommentCount(t *testing.T) {
	core.Init()

	// todo need test
	wg.Add(1)
	row := storage.Tencent.GetDetailUrl()[0]
	tencentCommentCount(row.Url, row.ShowId)
}
