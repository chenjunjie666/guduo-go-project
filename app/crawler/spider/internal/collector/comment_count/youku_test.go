package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestYoukuCommentCount(t *testing.T) {
	core.Init()
	// todo need complete
	wg.Add(1)
	//row := storage.Youku.GetDetailUrl()[0]
	youkuCommentCount("https://v.youku.com/v_show/id_XODA0MTU0NjM2.html?s=a6e80bc0fc0911e38b3f", 0)
}
