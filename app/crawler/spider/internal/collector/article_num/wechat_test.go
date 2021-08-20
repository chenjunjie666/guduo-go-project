package article_num

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/pkg/util"
	"testing"
)

func TestWechatArticleNum(t *testing.T) {
	core.Init()

	wg.Add(1)
	ch.PushJob()
	u := fmt.Sprintf("https://weixin.sogou.com/weixin?type=2&query=%s", util.UrlEncode("司藤"))
	wechatArticleNum(u, 0)

	wg.Wait()
}

