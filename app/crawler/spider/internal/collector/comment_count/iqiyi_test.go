package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestIqiyiCommentCount(t *testing.T) {
	core.Init()

	// todo need test
	wg.Add(1)
	ch.PushJob()
	//row := storage.Iqiyi.GetDetailUrl()[0]
	//iqiyiCommentCount(row.Url, row.ShowId)
	iqiyiCommentCount("https://www.iqiyi.com/v_2f1pqm3y53c.html", 0)    // 电影
	// iqiyiCommentCount("https://www.iqiyi.com/v_wu9wftmpq4.html", 100)    // 电视剧
	// iqiyiCommentCount("https://www.iqiyi.com/v_1vqcunfvcuw.html", 100)   // 综艺
	// iqiyiCommentCount("https://www.iqiyi.com/v_1whaybd4dn4.html", 100)   // 动漫
}
