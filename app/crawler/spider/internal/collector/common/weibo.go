package common

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/crawler/spider/internal/util"
	"net/http"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var Weibo = &weibo{
	PlatformId: storage.Weibo.PlatformId,
	Host:       storage.Weibo.Host,
	ApiHosts: struct {
		DataHost string // 微博指数用的api host
	}{
		DataHost: "https://data.weibo.com",
	},
}

type weibo struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		DataHost string // 微博指数用的api host
	}
}

func (w weibo) MobileCollector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"微博",
		w.Host,
		w.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	// 设置 cookie
	//cookies := w.getCookies()

	// 设置 cookie
	//c.OnRequest(func(r *colly.Request) {
	//r.Headers.Set("Cookie", util.BuildCookie(cookies))
	//})

	extensions.RandomMobileUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                                  // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                        // 非utf-8字符集支持
	return c
}

// 微博采集器初始化
func (w weibo) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"微博",
		w.Host,
		w.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	// 设置 cookie
	cookies := w.getCookies()

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
	})

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (w weibo) CollectorWithLogin(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"微博",
		w.Host,
		w.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	// 设置 cookie
	cookies := w.getLoginCookie()

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
		r.Headers.Set("Host", "s.weibo.com")
		r.Headers.Set("Host", "s.weibo.com")
	})

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c

}

func (w weibo) getCookies() []http.Cookie {
	cookies := []http.Cookie{
		//{
		//	Name:  "SUB",
		//	Value: "_2A25NXHamDeRhGeRN6lEW8CzLzTyIHXVuKO9urDV8PUNbmtANLU3akW9NU6FLx1JxP7-9cIVgRPApC6aXXg95r4KT",
		//},
		//{
		//	Name:  "SUBP",
		//	Value: "0033WrSXqPxfM725Ws9jqgMF55529P9D9WWMhQwZX1i0.p8VbwkmVdAH5JpX5KMhUgL.Foz0eKeNehzNSo52dJLoI7Dy9riyIgxkIrWr",
		//},
		{
			Name:  "SUB",
			Value: "_2A25NlX-eDeRhGeVG6VsZ9yzJyTWIHXVu49ZWrDV8PUNbmtANLWHskW9NT6E4UI9F1vpn3QmfQeH4vGOihx4Hro0c",
		},
		{
			Name:  "SUBP",
			Value: "0033WrSXqPxfM725Ws9jqgMF55529P9D9WFAlbDgzudeZQf_wIpy64fr5JpX5KzhUgL.FoeReo.RS0zfeo.2dJLoI7yXIsHk-sHEIBtt",
		},
	}

	return cookies
}

func (w weibo) getLoginCookie() []http.Cookie {

	cookies := []http.Cookie{
		{
			Name:  "ALF",
			Value: "1651892697",
		},
		{
			Name:  "SSOLoginState",
			Value: "1620356697",
		},
		{
			Name:  "SUB",
			Value: "_2A25NkN4JDeRhGeNG7FYQ-S_LzTiIHXVu5EjBrDV8PUNbmtANLWLEkW9NSzCmlV5h7qeIQLImPc5KHnqUrFDzyaKh",
		},
		{
			Name:  "SUBP",
			Value: "0033WrSXqPxfM725Ws9jqgMF55529P9D9WFLdddN6i7z.WfsK-MP32ly5JpX5KzhUgL.Fo-RS0Bp1K2NSoB2dJLoIE2LxKBLBonL1h.LxK-LBo.LBoxFi--Ni-iFiKnXe0zN1K-t",
		},
	}

	return cookies
}

func (w weibo) FetchIndicatorSearchResUrl(name string) string {
	c := w.MobileCollector("wid")
	req_data := map[string]string{
		"word": name,
	}

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	wid := ""
	c.OnResponse(func(r *colly.Response) {
		j := r.Body
		h, e := jsonparser.GetString(j, "html")
		if e != nil {
			log.Warning(fmt.Sprintf("未获取到搜索结果：%s，源数据：%s", e, j))
			return
		}

		reg := regexp.MustCompile(`wid="\d+"`)
		widTmp := reg.FindString(h)
		wid = strings.Trim(widTmp, "wid=\"")
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("referer", "https://data.weibo.com/index/ajax/newindex/searchword")
	})

	_ = c.Post("https://data.weibo.com/index/ajax/newindex/searchword", req_data)

	return fmt.Sprintf("https://data.weibo.com/index/newindex?visit_type=trend&wid=%s", wid)
}

func (w weibo) ParseWid(u string) string {
	reg := regexp.MustCompile(`wid=\d+`)

	widTmp := reg.FindString(u)

	wid := strings.Trim(widTmp, "wid=")

	return wid
}
