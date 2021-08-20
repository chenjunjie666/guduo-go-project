package common

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var Souhu = &souhu{
	PlatformId: storage.Souhu.PlatformId,
	Host:       storage.Souhu.Host,
	ApiHosts: struct {
		ApiHost      string
		PlHdHost     string
		DanmuApiHost string
	}{
		ApiHost:      "https://api.my.tv.sohu.com",
		PlHdHost:     "https://pl.hd.sohu.com",
		DanmuApiHost: "https://api.danmu.tv.sohu.com",
	},
}

type souhu struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		ApiHost      string
		PlHdHost     string
		DanmuApiHost string
	}
}

// 搜狐采集器初始化
func (s souhu) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"搜狐视频",
		s.Host,
		s.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持 - 最开始添加utf8支持的网站
	return c
}

// 搜狐的callback参数生成
func (s souhu) GenJQueryCallbackStr() string {
	ns := time.Now().UnixNano()       // 一个19位数字，直接用纳秒
	ms := time.Now().UnixNano() / 1e6 // 13位毫秒级时间戳
	cbNo := fmt.Sprintf("jQuery%d_%d", ns, ms)

	return cbNo
}

// 解析搜狐TV的vid
func (s souhu) ParseVid(u string) (string, string) {
	c := s.Collector("详情页获取VID")

	vid := ""
	html := ""

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		html = string(r.Body)

		reg := regexp.MustCompile(`vid="\d+"`)
		vid = reg.FindString(html)

		vid = strings.Trim(vid, `vid="`)

		if vid == "" {
			log.Warning("搜狐未能获取到vid，源连接为:", u, "返回的内容为：", html)
		}
	})
	c.Visit(u)

	return vid, html
}

func (s souhu) ParseVids(u string) []string {
	vids := make([]string, 0, 10)

	// 从详情页获取视频的vid
	// 搜狐视频的连接格式：
	apiHost := s.ApiHosts.PlHdHost
	playlistId := s.ParsePlayListId(u)
	o_playlistId := s.ParseOPlayListId(u)
	// 经过观察此参数是固定值
	pagesize := "999"
	cbNo := s.GenJQueryCallbackStr()

	// 格式化获取播放量的链接
	u = fmt.Sprintf("%s/videolist?playlistid=%s&o_playlistId=%s&callback=%s&pagesize=%s",
		apiHost,      // 网站host
		playlistId,   // playlistid 参数
		o_playlistId, // o_playlistId 参数
		cbNo,         // callback 参数
		pagesize,     // callback 参数
	)
	c := s.Collector("搜狐视频详情接口")

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	c.OnResponse(func(r *colly.Response) {
		// 返回值格式为 jsonp_xxxx_xxxx(json_str)
		// 所以需要将返回值转为 string
		// 然后去掉 "jsonp_xxxx_xxxx(" 以及 ")" 去掉，得到正确的 json 字符串
		ctxStr := string(r.Body)

		ctxStr = strings.TrimLeft(ctxStr, cbNo+"(")
		ctxStr = strings.TrimRight(ctxStr, ");")
		ctxByte := []byte(ctxStr)

		_, _ = jsonparser.ArrayEach(ctxByte, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

			vidTmp, e := jsonparser.GetString(value, "pageUrl")

			if e != nil {
				log.Warn(fmt.Sprintf("获取搜狐vid失败：%s, 连接为：%s, 原文是：%s", e, u, ctxStr))
				return
			}

			vid, html := s.ParseVid(vidTmp)
			if vid == "" {
				log.Warn(fmt.Sprintf("从URL获取搜狐vid失败：%s, 连接为：%s, 返回的原文: %s", vidTmp, u, html))
				return
			}

			vids = append(vids, vid)
		}, "videos")
	})

	_ = c.Visit(u)

	return vids
}

// 解析搜狐TV的playlistId 有时候也称为aid
func (s souhu) ParsePlayListId(u string) string {
	if regexp.MustCompile(`film\.sohu\.com`).FindString(u) != "" {
		uArr := strings.Split(u, "/")
		last := uArr[len(uArr) - 1]
		last = strings.Split(last, "?")[0]
		playlistId := strings.Split(last, ".")[0]
		return playlistId
	}

	c := s.Collector("详情页-ParsePlayListId")

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "tv.sohu.com")
	})

	playlistId := ""

	c.OnResponse(func(r *colly.Response) {
		html := string(r.Body)
		reg := regexp.MustCompile(`playlistId="\d+"`)
		playlistId = reg.FindString(html)
		//fmt.Println(html)

		playlistId = strings.Trim(playlistId, `playlistId="`)
	})

	_ = c.Visit(u)

	return playlistId
}

// 解析搜狐TV的o_playlistId
func (s souhu) ParseOPlayListId(u string) string {
	c := s.Collector("详情页-ParseOPlayListId")

	o_playlistId := ""

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	c.OnResponse(func(r *colly.Response) {
		html := string(r.Body)

		reg := regexp.MustCompile(`o_playlistId="\d+"`)
		o_playlistId = reg.FindString(html)

		o_playlistId = strings.Trim(o_playlistId, `o_playlistId="`)
	})

	_ = c.Visit(u)

	return o_playlistId
}
