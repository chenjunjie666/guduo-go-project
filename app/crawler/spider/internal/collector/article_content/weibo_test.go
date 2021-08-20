package article_content

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestWeiboArticleContent(t *testing.T) {
	core.Init()

	wg.Add(1)
	//url := storage.Weibo.GetFetchArticleUrl()[0]
	//weiboArticleContent(url.Url, url.ID)


	weiboArticleContent("https://s.weibo.com/weibo/%E5%B0%8F%E8%88%8D%E5%BE%97?topnav=1&wvr=6&b=1#_loginLayer_1620356647940", 19202)
}
