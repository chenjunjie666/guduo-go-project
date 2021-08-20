package danmaku

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core"
	"strings"
	"testing"
)

func TestYoukuDanmaku(t *testing.T) {
	core.Init()

	wg.Add(1)
	ch.PushJob()
	// 综艺
	//u := "https://v.youku.com/v_show/id_XNTEzOTI1NTk0OA==.html?spm=a2ha1.14919748_WEBZY_JINGXUAN.drawer5.d_zj1_2&s=cbcb9fd84b0e4cdebf60&scm=20140719.apircmd.5596.show_cbcb9fd84b0e4cdebf60"

	// 电视剧
	u := "https://v.youku.com/v_show/id_XNTAzNDM2MDY5Ng==.html?spm=a2h0c.8166622.PhoneSokuProgram_1.dselectbutton_1&showid=1e61efbfbdefbfbd04ef"

	// 电影
	//u := "https://v.youku.com/v_show/id_XNTEyMTU4NDE4MA==.html?spm=a2ha1.14919748_WEBMOVIE_JINGXUAN.drawer5.d_zj1_2&s=20efbfbd52055c6211ef&scm=20140719.apircmd.4424.show_20efbfbd52055c6211ef&s=20efbfbd52055c6211ef"

	// 动漫
	//u := "https://v.youku.com/v_show/id_XNDM3ODM0MDM2NA==.html?spm=a2ha1.14919748_WEBCOMIC_JINGXUAN.drawer4.d_zj1_2&s=d6bc38efbfbdefbfbdef&scm=20140719.manual.4392.show_d6bc38efbfbdefbfbdef&s=d6bc38efbfbdefbfbdef"

	//u ="https://list.youku.com/show/id_zffe36e46bc6c46e795b3.html"
	u = filterUrl(u)
	youkuDanmakuContent(u, 0)
}


func filterUrl(u string) string {
	if !strings.Contains(u, "list") {
		return u
	}
	cInfo := &core.CollectorInfo{
		"优酷",
		common.Youku.Host,
		common.Youku.PlatformId,
		"优酷视频详情页",
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
	c.OnHTML(".p-thumb a", func(ele *colly.HTMLElement) {
		link = ele.Attr("href")
		link = "https:" + link
	})

	c.Visit(u)
	return link
}