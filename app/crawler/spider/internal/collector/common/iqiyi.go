package common

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/internal/model_scrawler/show_detail_model"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var Iqiyi = &iqiyi{
	PlatformId: storage.Iqiyi.PlatformId,
	Host:       storage.Iqiyi.Host,
	ApiHosts: struct {
		SnsHost    string
		AlbumHost  string
		BulletHost string
	}{
		SnsHost:    "https://sns-comment.iqiyi.com",
		AlbumHost:  "https://pcw-api.iqiyi.com",
		BulletHost: "https://cmts.iqiyi.com",
	},
}

type iqiyi struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		SnsHost    string
		AlbumHost  string
		BulletHost string
	}
}

// 爱奇艺采集器初始化
func (i iqiyi) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"爱奇艺",
		i.Host,
		i.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	//c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (i iqiyi) GenJsonpCallbackStr() string {
	// 生成 callback 参数
	ms := time.Now().UnixNano() / 1e6
	rs := time.Now().Unix() / 1e5 // 五位数字，这里直接用10位时间戳前五位
	cb := fmt.Sprintf("jsonp_%d_%d", ms, rs)

	return cb
}

func (i iqiyi) ParseVid(u string) string {
	c := i.Collector("详情页")

	vid := ""

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	// 从详情页提取 vid，有些时候叫 content id
	c.OnResponse(func(r *colly.Response) {
		html := r.Body

		reg := regexp.MustCompile(`'albumid'] = "\d+"`)
		tvidTmp := reg.Find(html)

		reg = regexp.MustCompile(`\d+`)
		vidByte := reg.Find(tvidTmp)

		if string(vidByte) == "0" {
			reg = regexp.MustCompile(`'tvid'] = "\d+"`)
			tvidTmp = reg.Find(html)

			reg = regexp.MustCompile(`\d+`)
			vidByte = reg.Find(tvidTmp)
		}

		vid = string(vidByte)
	})

	_ = c.Visit(u)

	return vid
}

func (i iqiyi) ParseVids(u string) []string {

	c := i.Collector("爱奇艺VID")

	vids := make([]string, 0, 10)
	aid := i.ParseAID(u)

	if aid == "" {
		log.Warn(fmt.Sprintf("爱奇艺，未获取到aid：%s", u))
		return vids
	}

	breakFlag := false
	page := 1
	size := "30"

	//return vids

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	// 从详情页提取 vid，有些时候叫 content id
	c.OnResponse(func(r *colly.Response) {
		jTmp := string(r.Body)

		reg := regexp.MustCompile(`\);}catch\(e\).*?$`)
		jTmp = reg.ReplaceAllString(jTmp, "")

		reg = regexp.MustCompile(`try.*?\(`)
		jTmp = reg.ReplaceAllString(jTmp, "")

		_, e := jsonparser.GetString([]byte(jTmp), "data", "epsodelist", "[0]", "vid")
		if e != nil {
			breakFlag = true
			return
		}

		_, _ = jsonparser.ArrayEach([]byte(jTmp), func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			vid, e2 := jsonparser.GetInt(v, "tvId")
			if e2 != nil {
				log.Warn(fmt.Sprintf("解析vid失败：%s", e))
				return
			}
			vids = append(vids, strconv.FormatInt(vid, 10))
		}, "data", "epsodelist")
	})

	for {
		callback := i.GenJsonpCallbackStr()
		apiUrl := fmt.Sprintf("%s/albums/album/avlistinfo?aid=%s&page=%d&size=%s&callback=%s",
			i.ApiHosts.AlbumHost,
			aid,
			page,
			size,
			callback,
		)
		_ = c.Visit(apiUrl)

		if breakFlag == true{
			break
		}
		page++
	}
	// 如果没有获取到vids那可能就是综艺之类的
	if len(vids) == 0 {
		c2 := i.Collector("爱奇艺VID-方式2")
		c2.OnError(func(r *colly.Response, e error) {
			c2.Retry(r, e)
		})
		c2.OnHTML(".side-content:first-child", func(ele *colly.HTMLElement) {
			ele.ForEach(".playing-icon", func(i int, ele *colly.HTMLElement) {
				attr := ele.Attr("v-show")
				vid := strings.Trim(attr, "(== itemLeft.tvId)")
				if vid != "" {
					vids = append(vids, vid)
				}
			})
		})
		c2.Visit(u)
	}

	return vids
}

func (i iqiyi) ParseAID(u string) string {
	c := i.Collector("详情页")

	aid := ""

	c.OnError(func(r *colly.Response, e error) {
		if strings.Contains(r.Request.URL.String(), "iqiyi.com/v_") && r.StatusCode == 404 {
			show_detail_model.ReportErrorUrl(r.Request.URL.String())
			return
		}
		c.Retry(r, e)
	})

	// 从详情页提取 aid
	c.OnResponse(func(r *colly.Response) {
		html := r.Body

		reg := regexp.MustCompile(`param\['albumid'.*?"\d+"`)
		aidTmp := reg.Find(html)

		reg = regexp.MustCompile(`\d+`)
		aidByte := reg.Find(aidTmp)

		aid = string(aidByte)
	})

	_ = c.Visit(u)

	return aid
}



func (i iqiyi) ParseLength(u string) map[string]int64 {

	c := i.Collector("获取时长")

	length := make(map[string]int64)
	aid := i.ParseAID(u)

	if aid == "" {
		log.Warn(fmt.Sprintf("爱奇艺，未获取到aid：%s", u))
		return make(map[string]int64)
	}

	breakFlag := false
	page := 1
	size := "30"

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	// 从详情页提取 vid，有些时候叫 content id
	c.OnResponse(func(r *colly.Response) {
		jTmp := string(r.Body)

		reg := regexp.MustCompile(`\);}catch\(e\).*?$`)
		jTmp = reg.ReplaceAllString(jTmp, "")

		reg = regexp.MustCompile(`try.*?\(`)
		jTmp = reg.ReplaceAllString(jTmp, "")

		_, e := jsonparser.GetString([]byte(jTmp), "data", "epsodelist", "[0]", "vid")
		if e != nil {
			breakFlag = true
			return
		}

		_, _ = jsonparser.ArrayEach([]byte(jTmp), func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			vid, e2 := jsonparser.GetInt(v, "tvId")
			if e2 != nil {
				log.Warn(fmt.Sprintf("解析vid失败：%s", e))
				return
			}

			len2, _ := jsonparser.GetString(v, "duration")
			lenArr := strings.Split(len2, ":")
			llArr := make([]int64, 0, 3)
			for _, l := range lenArr {
				ll, _ := strconv.ParseInt(l, 10, 64)
				llArr = append(llArr, ll)
			}

			second := int64(0)
			if len(llArr) == 3 {
				second = llArr[0] * 3600 + llArr[1] * 60 + llArr[2]
			}else if len(llArr) == 2 {
				second = llArr[0] * 60 + llArr[1]
			}else if len(llArr) == 1 {
				second = llArr[0]
			}else{
				second = 4 * 3600
			}

			lTmp := second / 300 + 1

			vidStr := strconv.FormatInt(vid, 10)
			length[vidStr] =lTmp
		}, "data", "epsodelist")
	})

	for {
		callback := i.GenJsonpCallbackStr()
		apiUrl := fmt.Sprintf("%s/albums/album/avlistinfo?aid=%s&page=%d&size=%s&callback=%s",
			i.ApiHosts.AlbumHost,
			aid,
			page,
			size,
			callback,
		)
		_ = c.Visit(apiUrl)

		if breakFlag == true{
			break
		}
		page++
	}

	if len(length) == 0 {
		c2 := i.Collector("爱奇艺VID-方式2")
		c2.OnError(func(r *colly.Response, e error) {
			c2.Retry(r, e)
		})
		c2.OnHTML(".side-content > div:first-child", func(ele *colly.HTMLElement) {
			attr := ele.Attr(":initialized-data")
			_, _ = jsonparser.ArrayEach([]byte(attr), func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
				vid, _ := jsonparser.GetInt(v, "tvId")
				dur, _ := jsonparser.GetInt(v, "duration")


				sec := dur / 300 + 1


				vidStr := strconv.FormatInt(vid, 10)
				length[vidStr] = sec
			})
		})
		c2.Visit(u)
	}

	return length
}

