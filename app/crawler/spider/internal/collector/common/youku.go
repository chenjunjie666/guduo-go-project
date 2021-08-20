package common

import (
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"guduo/app/crawler/spider/internal/core"
	youku2 "guduo/app/crawler/spider/internal/lib/youku"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/crawler/spider/internal/util"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var Youku = &youku{
	PlatformId: storage.Youku.PlatformId,
	Host:       storage.Youku.Host,
	AppKey:     youku2.GetAppKey(),
}

type youku struct {
	PlatformId uint64
	Host       string
	AppKey     int
}

// 优酷采集器初始化
func (y youku) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"优酷",
		y.Host,
		y.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (y youku) CollectorWithToken(mod string) *core.CollectorObj {
	MH5Tk, MH5TkEnc := youku2.GetToken()

	cInfo := &core.CollectorInfo{
		"优酷",
		y.Host,
		y.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	// 设置 cookie
	cookies := y.getCookies(MH5Tk, MH5TkEnc)

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
		r.Headers.Add("Content-type", "application/x-www-form-urlencoded")
		r.Headers.Add("Refer", "https://v.youku.com/")
	})

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (y youku) getCookies(MH5Tk, MH5TkEnc string) []http.Cookie {
	cookies := []http.Cookie{
		{
			Name:  "_m_h5_tk",
			Value: MH5Tk,
		},
		{
			Name:  "_m_h5_tk_enc",
			Value: MH5TkEnc,
		},
	}

	return cookies
}

// 获取sid
func (y youku) ParseSid(u string) string {
	sid := ""
	sidTmpArr := strings.Split(u, "?")
	if len(sidTmpArr) == 2 {
		sidTmp := sidTmpArr[1]
		sidTmpArr = strings.Split(sidTmp, "&")
		for _, v := range sidTmpArr {
			sidTmp = regexp.MustCompile(`^s(howId)?=`).FindString(v)
			if sidTmp != "" {
				sid = regexp.MustCompile(`^s(howId)?=`).ReplaceAllString(sidTmp, `^s(howId)?=`)
			}
		}
	}
	if sid == "" {
		c := y.Collector("详情页")
		c.OnError(func(r *colly.Response, e error) {
			c.Retry(r, e)
		})
		c.OnResponse(func(r *colly.Response) {
			body := string(r.Body)
			reg := regexp.MustCompile(`showid_en:\s?'.*?'`)
			sidTmp := reg.FindString(body)
			sidTmp = strings.Trim(sidTmp, "showid_en: ")
			sid = strings.Trim(sidTmp, "'") // 获取ShowId
		})
		_ = c.Visit(u)
	}

	return sid
}

// 获取vid
func (y youku) ParseVid(u string) string {
	vid := ""
	vidTmp := strings.Split(u, "?")[0]
	vidTmpArr := strings.Split(vidTmp, "/")
	for _, v := range vidTmpArr {
		if strings.Contains(v, ".html") {
			vidTmp = strings.TrimRight(v, ".html")
			vid = strings.TrimLeft(vidTmp, "id_") // 获取videoId
			break
		}
	}

	return vid
}

// 获取分集的VID
func (y youku) ParseVids(u string) []string {
	c := y.Collector("详情页")
	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	sid := y.ParseSid(u)
	vid := y.ParseVid(u)
	vids := make([]string, 0, 10)
	c.OnResponse(func(r *colly.Response) {
		body := string(r.Body)

		reg := regexp.MustCompile(`episodeLast: '\d+'`)
		epNumTmp := reg.FindString(body)
		epNumTmp = strings.Trim(epNumTmp, `episodeLast: `)
		epNumTmp = strings.Trim(epNumTmp, `'`) // 获取utdid
		epNum, _ := strconv.ParseInt(epNumTmp, 10, 64)
		epNums := int(epNum)

		if epNum > 10000 {
			return
		}
		page := 1
		for {
			apiUrl := youku2.GetEpisodeUrl(vid, sid, page)
			vidRow := y.parseVids(apiUrl)
			if len(vidRow) == 0 {
				break
			}
			vids = append(vids, vidRow...)
			if page * 30 >= epNums{
				break
			}
			page++
		}
	})

	vids2 := make([]string, 0, 10)
	c.OnHTML(".listbox", func(ele *colly.HTMLElement) {
		if strings.Contains(ele.Text, "VIP专享") {
			ele.ForEach(".anthology-content-scroll>div>.anthology-content:first-child .pic-text-item>a", func(i int, ele *colly.HTMLElement) {
				url := ele.Attr("href")
				vidRow := y.ParseVid(url)
				vids2 = append(vids2, vidRow)
			})
		}


	})

	_ = c.Visit(u)

	if len(vids2) > 0 {
		return vids2
	}

	if len(vids) > 0 {
		return vids
	}

	return []string{vid}
}

// 获取分集vid的子方法，当非综艺时使用
func (y youku) parseVids(u string) []string {
	c := y.CollectorWithToken("分集连接")
	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	vids := make([]string, 0, 10)
	c.OnResponse(func(r *colly.Response) {
		ret, _ := jsonparser.GetString(r.Body, "ret", "[0]")
		if !strings.Contains(ret, "SUCCESS::调用成功") {
			c.Retry(r, nil)
			return
		}
		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			vid, _ := jsonparser.GetString(v, "data", "action", "value")
			if vid != "" {
				vids = append(vids, vid)
			}

		}, "data", "2019030100", "data", "nodes")
	})

	_ = c.Visit(u)
	return vids
}

func (y youku) ParseLength(u string) int {
	c := y.Collector("详情页")
	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	minis := 0
	c.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		reg := regexp.MustCompile(`seconds:\s?'.*?'`)
		secTmp := reg.FindString(body)
		secTmp = strings.Trim(secTmp, `seconds: `)
		secTmp = strings.Trim(secTmp, `'`) // 获取utdid
		secIntTmp, _ := strconv.ParseFloat(secTmp, 32)
		mins := secIntTmp / 60
		minis = int(mins)
	})

	_ = c.Visit(u)

	return minis
}

//func (y youku) ParseVids(u string) []string {
//	j := y.ParseHtmlInitDataJson(u)
//
//}

// 解析原始html文档中的视频信息
// 这些信息被存储在了HTML中的一个json字符串中
func (y youku) ParseHtmlInitDataJson(u string) []byte {
	json := make([]byte, 0, 10)

	c := y.Collector("详情页")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		html := string(r.Body)
		reg := regexp.MustCompile(`__INITIAL_DATA__.*?;</script>`)

		flagStr := reg.FindString(html)

		j := strings.Trim(flagStr, "__INITIAL_DATA__ =;</script>")

		json = []byte(j)
	})

	_ = c.Visit(u)

	return json
}



















//type youkuDetailPageProp struct {
//	ShowId  string
//	VideoId string
//	Minis   int    // 视频分钟数
//	Cna     string // 分集接口需要的参数：utdid
//	EpNum   int    // 总集数
//}
//// 暂时不用
//func (y youku) ParseDetailPage(u string) *youkuDetailPageProp {
//	dpp := &youkuDetailPageProp{}
//
//	vidTmp := strings.Split(u, "?")[0]
//	vidTmpArr := strings.Split(vidTmp, "/")
//	for _, v := range vidTmpArr {
//		if strings.Contains(v, ".html") {
//			vidTmp = strings.TrimRight(v, ".html")
//			dpp.VideoId = strings.TrimLeft(vidTmp, "id_") // 获取videoId
//			break
//		}
//	}
//
//	c := y.Collector("详情页")
//	c.OnError(func(r *colly.Response, e error) {
//		c.Retry(r, e)
//	})
//
//	c.OnResponse(func(r *colly.Response) {
//		body := string(r.Body)
//		reg := regexp.MustCompile(`showid_en:\s?'.*?'`)
//		sidTmp := reg.FindString(body)
//		sidTmp = strings.Trim(sidTmp, "showid_en: ")
//		dpp.ShowId = strings.Trim(sidTmp, "'") // 获取ShowId
//
//		reg = regexp.MustCompile(`"cna":".*?"`)
//		cnaTmp := reg.FindString(body)
//		cnaTmp = strings.Trim(cnaTmp, `"cna`)
//		dpp.Cna = strings.Trim(cnaTmp, `":`) // 获取utdid
//
//		reg = regexp.MustCompile(`seconds:\s?'.*?'`)
//		secTmp := reg.FindString(body)
//		secTmp = strings.Trim(secTmp, `seconds: `)
//		secTmp = strings.Trim(secTmp, `'`) // 获取utdid
//		secIntTmp, e := strconv.ParseFloat(secTmp, 10)
//		if e == nil {
//			mins := secIntTmp / 60
//			dpp.Minis = int(mins) + 1
//		}
//
//		reg = regexp.MustCompile(`episodeLast: '\d+'`)
//		epNumTmp := reg.FindString(body)
//		epNumTmp = strings.Trim(epNumTmp, `episodeLast: `)
//		epNumTmp = strings.Trim(epNumTmp, `'`) // 获取utdid
//		epNum, e := strconv.ParseInt(epNumTmp, 10, 64)
//		if e == nil {
//			dpp.EpNum = int(epNum)
//		}
//	})
//
//	// 搜索 .anthology-content-scroll 这个class 下有没有 VIP专享 关键字，有就直接从这里找svid，否则从接口获取（集数类型）
//	// vip专享会返回对应的剧集连接
//	c.OnHTML(".anthology-content-scroll>div>.anthology-content:first-child", func(ele *colly.HTMLElement) {
//		ele.ForEach(".pic-text-item>a", func(i int, ele *colly.HTMLElement) {
//			url := ele.Attr("href")
//			fmt.Println(url)
//		})
//	})
//
//	_ = c.Visit(u)
//
//	return dpp
//}

func (y youku) CheckIsYoukuPlayPage(u string) bool {
	if strings.Contains(u, "youku.com/v_show") {
		return true
	}

	return false
}