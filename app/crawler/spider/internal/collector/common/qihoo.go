package common

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/crawler/spider/internal/util"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var Qihoo = &qihoo{
	PlatformId: storage.Qihoo.PlatformId,
	Host:       storage.Qihoo.Host,
	ApiHosts: struct {
		IndexHost string
	}{
		IndexHost: "https://trends.so.com",
	},
}

type qihoo struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		IndexHost string
	}
}

// 360采集器初始化
func (q qihoo) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"360",
		q.Host,
		q.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	// 设置 cookie
	cookies := q.getCookies()

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
	})

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (q qihoo) getCookies() []http.Cookie {
	cookies := []http.Cookie{
		{
			Name:  "T",
			Value: "s%3Df2dffefec56de59fd55355400d6a05ca%26t%3D1621441348%26lm%3D0-1%26lf%3D2%26sk%3Dea27809ecbc40ccffe034dcc75217bdb%26mt%3D1621441348%26rc%3D%26v%3D2.0%26a%3D1",
		}, {
			Name:  "Q",
			Value: "u%3D360H2742777278%26n%3D%26le%3D%26m%3DZGp3WGWOWGWOWGWOWGWOWGWOAmNk%26qid%3D2742777278%26im%3D1_t0105d6cf9b508f72c8%26src%3Dpcw_360index%26t%3D1",
		},
	}

	return cookies
}

// 解析360指标的搜索关键词
func (q qihoo) ParseWords(u string) string {
	words := ""
	reg := regexp.MustCompile(`(\?|&)query=.*?(&|$)`)
	words = reg.FindString(u)
	words = strings.Trim(words, `?query=&`)

	return words
}

// 解析360指标的搜索开始时间
func (q qihoo) ParseStartDate() string {
	now := time.Now()
	year := now.Year()   //年
	month := now.Month() //月
	day := now.Day()     //日
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local).AddDate(0, 0, -8)
	startDateStr := startDate.Format("20060102")

	return startDateStr
}

// 解析360指标的搜索开始时间
func (q qihoo) ParseEndDate() string {
	now := time.Now()
	year := now.Year()   //年
	month := now.Month() //月
	day := now.Day()     //日
	endDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	endDateStr := endDate.Format("20060102")

	return endDateStr
}
