package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestMangoCommentCount(t *testing.T) {
	core.Init()
	// todo need test
	wg.Add(1)
	//row := storage.Mango.GetDetailUrl()[0]


	mangoCommentCount("https://www.mgtv.com/b/364827/11598187.html", 0)
}
