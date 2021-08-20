package common

import (
	"bytes"
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2/extensions"
)

var Mango = &mango{
	PlatformId: storage.Mango.PlatformId,
	Host:       storage.Mango.Host,
	ApiHosts: struct {
		CommentHost      string
		VcHost           string
		PcWebHost        string
		BulletHost       string
		GalaxyBxHost     string
		GalaxyPersonHost string
	}{
		CommentHost:      "https://comment.mgtv.com",
		VcHost:           "https://vc.mgtv.com",
		PcWebHost:        "https://pcweb.api.mgtv.com",
		BulletHost:       "https://bullet-ali.hitv.com",
		GalaxyBxHost:     "https://galaxy.bz.mgtv.com",
		GalaxyPersonHost: "http://galaxy.person.mgtv.com",
	},
}

type mango struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		CommentHost      string
		VcHost           string
		PcWebHost        string
		BulletHost       string
		GalaxyBxHost     string
		GalaxyPersonHost string
	}
}

// 芒果TV采集器初始化
func (m mango) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"芒果TV",
		m.Host,
		m.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	return c
}

// 根据详情页链接，提取其中的 video id
func (m mango) ParseVid(u string) string {
	uParse := strings.Split(u, "/")
	uEnd := uParse[len(uParse)-1]
	uParse = strings.Split(uEnd, ".")

	vid := uParse[0]

	return vid
}

// 根据详情页链接，提取其中的 Category id
func (m mango) ParseCid(u string) string {
	uParse := strings.Split(u, "/")
	if len(uParse) < 2 {
		log.Warn(fmt.Sprintf("芒果TV，无法解析url的cid, %s", u))
		return ""
	}
	cid := uParse[len(uParse)-2]

	return cid
}

func (m mango) ParseVids(u string) []string {
	breakFlag := false

	_support := "10000000"
	version := "5.5.35"
	videoId := m.ParseVid(u)
	page := 0
	size := "30"
	callback := m.GenJsonpCallbackStr()

	vids := make([]string, 0, 10)
	accessCount := 0

	c := m.Collector("芒果TV vid")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		jTmp := string(r.Body)
		j := strings.Trim(jTmp, "jsonp_1234567890();")

		resCode, _ := jsonparser.GetInt([]byte(j), "code")
		if resCode != 200 {
			breakFlag = true
			log.Warning("获取分集列表失败，内容：", string(j))
			return
		}

		totalPage, _ := jsonparser.GetInt([]byte(j), "data", "total_page")
		curPage, _ := jsonparser.GetInt([]byte(j), "data", "current_page")

		_, _ = jsonparser.ArrayEach([]byte(j), func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			vid, e := jsonparser.GetString(v, "video_id")
			if e != nil {
				log.Warn(fmt.Sprintf("获取vid失败：%s", e))
				return
			}

			vids = append(vids, vid)
			accessCount++
		}, "data", "list")

		if curPage == totalPage {
			breakFlag = true
			return
		}
	})

	for {
		apiUrl := fmt.Sprintf("%s/episode/list?_support=%s&version=%s&video_id=%s&page=%d&size=%s&callback=%s",
			m.ApiHosts.PcWebHost,
			_support,
			version,
			videoId,
			page,
			size,
			callback,
		)
		_ = c.Visit(apiUrl)

		if breakFlag {
			break
		}
		page++
	}

	return vids
}

// 生成芒果TV的 jsonp_callback 参数的值
func (m mango) GenJsonpCallbackStr() string {
	ms := time.Now().UnixNano() / 1e6
	rs := time.Now().Unix() / 1e5 // 五位数字，这里直接用10位时间戳前五位
	cbNo := fmt.Sprintf("jsonp_%d_%d", ms, rs)

	return cbNo
}

// 解析的搜索时间
func (m mango) ParseCurrentDate() string {
	now := time.Now()
	currentDateStr := now.Format("2006/01/02")
	return currentDateStr
}

func (m mango) ParseVideoLength(vid string, cid string) int64 {
	c := m.Collector("单集片长")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})
	length := int64(0)
	findFlag := false
	_support := 10000000
	c.OnResponse(func(r *colly.Response) {
		findFlag = true
		body := bytes.Trim(r.Body, "jsonp_1234567890();")
		lengthTmp, _ := jsonparser.GetString(body, "data", "info", "time")
		lengthParse := strings.Split(lengthTmp, ":")
		hour := int64(0)
		second := int64(0)
		min := int64(0)
		if len(lengthParse) == 3 {
			second, _ = strconv.ParseInt(lengthParse[2], 10, 64)
			min, _ = strconv.ParseInt(lengthParse[1], 10, 64)
			hour, _ = strconv.ParseInt(lengthParse[0], 10, 64)
		} else if len(lengthParse) == 2 {
			second, _ = strconv.ParseInt(lengthParse[1], 10, 64)
			min, _ = strconv.ParseInt(lengthParse[0], 10, 64)
		}else if len(lengthParse) == 0{
			findFlag = false
		}else {
			second = 60
		}
		length = hour*3600 + min*60 + second

		// 最大时长为4小时
		if length > 3600 * 4 {
			length = 3600 * 4
		}

	})
	url := fmt.Sprintf("%s/video/info?vid=%s&cid=%s&_support=%d",
		m.ApiHosts.PcWebHost,
		vid,
		cid,
		_support,
	)
	_ = c.Visit(url)

	if !findFlag {
		log.Warn(fmt.Sprintf("未找到片长"))
		return 0
	}
	return length
}



func (m mango) ParseLen(u string) map[string]int64 {
	breakFlag := false
	_support := "10000000"
	version := "5.5.35"
	videoId := m.ParseVid(u)
	page := 0
	size := "30"
	callback := m.GenJsonpCallbackStr()

	ret := make(map[string]int64)

	c := m.Collector("芒果TV，获取视频时长")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		jTmp := string(r.Body)

		j := strings.Trim(jTmp, callback+")")

		totalPage, _ := jsonparser.GetInt([]byte(j), "data", "total_page")
		curPage, _ := jsonparser.GetInt([]byte(j), "data", "current_page")
		if curPage == totalPage {
			breakFlag = true
		}

		_, _ = jsonparser.ArrayEach([]byte(j), func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			vid, e := jsonparser.GetString(v, "video_id")
			if e != nil {
				log.Warn(fmt.Sprintf("获取vid失败：%s", e))
				return
			}

			len2, _ := jsonparser.GetString(v, "time")
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

			ret[vid] = second
		}, "data", "list")
	})

	for {
		apiUrl := fmt.Sprintf("%s/episode/list?_support=%s&version=%s&video_id=%s&page=%d&size=%s&&callback=%s",
			m.ApiHosts.PcWebHost,
			_support,
			version,
			videoId,
			page,
			size,
			callback,
		)
		_ = c.Visit(apiUrl)
		if breakFlag {
			break
		}
		page++
	}

	return ret
}