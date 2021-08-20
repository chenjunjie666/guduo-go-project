package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestSouhuCommentCount(t *testing.T) {
	core.Init()
	// todo need test
	wg.Add(1)
	row := storage.Souhu.GetDetailUrl()[0]
	souhuCommentCount(row.Url, row.ShowId)
}
