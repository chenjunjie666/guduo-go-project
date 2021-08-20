package common

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/crawler/spider/internal/util"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
)

var Bilibili = &bilibili{
	PlatformId: storage.Bilibili.PlatformId,
	Host:       storage.Bilibili.Host,
	ApiHosts: struct {
		ApiHost string
	}{
		ApiHost: "https://api.bilibili.com",
	},
}

type bilibili struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		ApiHost string
	}
}

// 哔哩哔哩采集器初始化
func (b bilibili) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"哔哩哔哩",
		b.Host,
		b.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫
	// 设置 cookie
	cookies := b.getCookies()

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
	})

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	//c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (b bilibili) getCookies() []http.Cookie {
	cookies := []http.Cookie{
		{
			Name:  "SESSDATA",
			Value: "33fcfb41%2C1633657869%2Cc3e2f%2A41",
		},
	}

	return cookies
}

func (b bilibili) GenJsonpCallbackStr() string {
	// 生成 callback 参数
	ms := time.Now().UnixNano() / 1e6
	rs := time.Now().Unix() / 1e5 // 五位数字，这里直接用10位时间戳前五位
	cb := fmt.Sprintf("jsonp_%d_%d", ms, rs)

	return cb
}

// 解析bilibili的搜索时间
func (b bilibili) ParseDate() string {
	now := time.Now()
	year := now.Year()   //年
	month := now.Month() //月
	day := now.Day()     //日
	Date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	DateStr := Date.Format("2006-01-02")

	return DateStr
}

// 根据b站的视频详情页链接解析ep_id
func (b bilibili) ParseEpid(u string) string {
	reg := regexp.MustCompile(`ep\d+\?`)
	epidTmp := reg.FindString(u)

	epid := strings.Trim(epidTmp, "ep?")
	return epid
}

// 解析获得所有集数的EPID
func (b bilibili) ParseEpids(u string) []string {
	c := b.Collector("详情页")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	season_id := ""
	epids := make([]string, 0, 10)
	// b站将所有的集数ID epid 放在了html里的一个json字符串中
	// 通过寻找特征值，找出这段json，在执行json搜索，得到所有epid
	c.OnResponse(func(r *colly.Response) {
		htmlByte := r.Body

		// 找到特征值
		reg := regexp.MustCompile(`"season_id":\d+`)
		jTmp := reg.Find(htmlByte)

		// 去除特征值字符，只留下正确的season id
		season_id = strings.Trim(string(jTmp), "season_id\":")

	})
	_ = c.Visit(u)

	apiUrl := fmt.Sprintf("https://api.bilibili.com/pgc/web/season/section?season_id=%s&season_type=4", season_id)
	c = b.Collector("season_id")
	c.OnResponse(func(r *colly.Response) {
		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			epid, e := jsonparser.GetInt(v, "id")
			if e != nil {
				log.Warning(fmt.Sprintf("获取epid失败：%s，源数据：%s", e, v))
			}
			epids = append(epids, strconv.FormatInt(epid, 10))
		}, "result", "main_section", "episodes")
	})
	_ = c.Visit(apiUrl)
	return epids
}

// 根据b站的视频详情页链接解析bvid
func (b bilibili) ParseBvid(u string) string {
	reg := regexp.MustCompile(`/BV\w+?(/|\?|$)`)
	bvidTmp := reg.FindString(u)
	bvid := strings.Trim(bvidTmp, `?/`)
	return bvid
}

// 根据b站的bvid获取aid 有时也叫pid
func (b bilibili) ParseAid(bvid string) string {
	apiUrl := fmt.Sprintf("%s/x/web-interface/view?bvid=%s",
		b.ApiHosts.ApiHost,
		bvid,
	)

	c := b.Collector("aid")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	aid := ""
	c.OnResponse(func(r *colly.Response) {
		findFlag = true

		html := r.Body
		aidInt, _ := jsonparser.GetInt(html, "data", "aid")
		aid = strconv.FormatInt(aidInt, 10)
	})

	_ = c.Visit(apiUrl)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到aid", apiUrl))
	}

	return aid
}

// 根据b站的bvid获取Cid(默认获取第一个p) 有时也叫oid
func (b bilibili) ParseCid(bvid string) string {
	apiUrl := fmt.Sprintf("%s/x/web-interface/view?bvid=%s",
		b.ApiHosts.ApiHost,
		bvid,
	)

	c := b.Collector("bilibili cid")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	cid := ""
	c.OnResponse(func(r *colly.Response) {
		findFlag = true

		html := r.Body
		cidInt, _ := jsonparser.GetInt(html, "data", "cid")
		cid = strconv.FormatInt(cidInt, 10)
	})

	_ = c.Visit(apiUrl)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到cid", apiUrl))
	}

	return cid
}

// bilibili每个视频可能有多个p(理解为一种剧集)
// 根据b站的bvid获取cid列表
//func (b bilibili) ParseCids(u string) []string {
//	var cids []string
//	cb := b.GenJsonpCallbackStr()
//	bvid := u
//	apiUrl := fmt.Sprintf("%s/x/player/pagelist?bvid=%s&jsonp=%s",
//		b.ApiHosts.ApiHost,
//		bvid,
//		cb,
//	)
//
//	c := b.Collector("cid")
//
//	c.OnError(func(r *colly.Response, e error) {
//		c.Retry(r, e)
//	})
//
//	findFlag := false
//
//	c.OnResponse(func(r *colly.Response) {
//		findFlag = true
//
//		html := r.Body
//
//		_, _ = jsonparser.ArrayEach(html, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
//			cid, _ := jsonparser.GetInt(value, "cid")
//			cids = append(cids, strconv.FormatInt(cid, 10))
//		}, "data")
//
//	})
//
//	_ = c.Visit(apiUrl)
//
//	if !findFlag {
//		log.Warn(fmt.Sprintf("链接:%s，未找到cid列表", u))
//	}
//
//	return cids
//}

func (b bilibili) ParseCids(u string) []string {
	c := b.Collector("详情页")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	season_id := ""
	cids := make([]string, 0, 10)
	// b站将所有的集数ID epid 放在了html里的一个json字符串中
	// 通过寻找特征值，找出这段json，在执行json搜索，得到所有epid
	c.OnResponse(func(r *colly.Response) {
		htmlByte := r.Body

		// 找到特征值
		reg := regexp.MustCompile(`"season_id":\d+`)
		jTmp := reg.Find(htmlByte)

		// 去除特征值字符，只留下正确的season id
		season_id = strings.Trim(string(jTmp), "season_id\":")
	})
	_ = c.Visit(u)

	apiUrl := fmt.Sprintf("https://api.bilibili.com/pgc/web/season/section?season_id=%s&season_type=4", season_id)
	c = b.Collector("season_id")
	c.OnResponse(func(r *colly.Response) {
		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			cid, e := jsonparser.GetInt(v, "cid")
			if e != nil {
				log.Warning(fmt.Sprintf("获取epid失败：%s，源数据：%s", e, v))
			}
			cids = append(cids, strconv.FormatInt(cid, 10))
		}, "result", "main_section", "episodes")
	})
	_ = c.Visit(apiUrl)
	return cids
}

// 根据b站的番剧详情页链接解析season_id
func (b bilibili) ParseSeasonId(u string) string {
	reg := regexp.MustCompile(`/ss\d+?(/|\?|$)`)
	seasonIdTmp := reg.FindString(u)
	seasonId := strings.Trim(seasonIdTmp, `?/`)
	seasonId = strings.TrimLeft(seasonId, `ss`)
	return seasonId
}
