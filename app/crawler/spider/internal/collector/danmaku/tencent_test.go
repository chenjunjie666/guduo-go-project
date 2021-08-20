package danmaku

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core"
	"strings"
	"testing"
)

func TestTencentDanmakuContent(t *testing.T) {
	core.Init()
	wg.Add(1)
	ch.PushJob()
	//url := storage.Tencent.GetDetailUrl()[0]
	// 数据库中的司藤的详情页
	//tencentHandle()
	//wg.Wait()
	//tencentDanmakuContent(filterUrlTencent("https://v.qq.com/detail/m/mzc00200vu38bvg.html"), 21908)

	//tencentDanmakuContent(filterUrlTencent("https://v.qq.com/detail/m/mzc00200vu38bvg.html"), 21908)
	tencentDanmakuContent("https://v.qq.com/x/cover/mzc00200yu62ksg/b0036vn707y.html", 99)

}



func filterUrlTencent(u string) string {
	if !strings.Contains(u, "detail"){
		return u
	}
	cInfo := &core.CollectorInfo{
		"腾讯视频",
		common.Tencent.Host,
		common.Tencent.PlatformId,
		"腾讯视频详情页",
	}
	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	link := ""
	c.OnHTML(".mod_episode .item:first-child a", func(ele *colly.HTMLElement) {
		link = ele.Attr("href")
	})

	c.Visit(u)
	return link
}