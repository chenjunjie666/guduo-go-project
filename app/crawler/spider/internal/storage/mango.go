package storage

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_scrawler/comment_count_model"
	"guduo/app/internal/model_scrawler/danmaku_count_model"
	"guduo/app/internal/model_scrawler/danmaku_model"
	"guduo/app/internal/model_scrawler/play_count_model"
	"guduo/app/internal/model_scrawler/show_detail_model"
	"strings"
	"sync"
	"time"
)

var Mango = &mango{
	PlatformId: constant.PlatformIdMango,
	Host:       "https://www.mgtv.com",
}

type mango struct {
	PlatformId uint64
	Host       string
}

// 获取芒果TV剧集的详情页链接
func (m *mango) GetDetailUrl() show_detail_model.DetailUrls {
	//showIds := show_model.GetActiveShows()
	//urls := show_detail_model.GetDetailUrl(m.PlatformId, showIds)
	urls := show_detail_model.GetDetailUrlNew(m.PlatformId)
	wg2 := &sync.WaitGroup{}
	for k, v := range urls {
		oUrl := v.Url
		wg2.Add(1)
		go func(k int, sid uint64) {
			urls[k].Url = m.filterUrl(oUrl)
			if urls[k].Url != oUrl {
				show_detail_model.SaveTrueUrl(urls[k].Url, sid, constant.PlatformIdMango)
			}
			wg2.Done()
		}(k, v.ShowId)
	}
	wg2.Wait()

	log.Info("芒果TV一共", len(urls), "部剧综")
	return urls
	//urls := []string{
	//	"https://www.mgtv.com/b/340159/9777172.html",
	//	"https://www.mgtv.com/b/337307/9552358.html",
	//}
	//return urls
}

func (m *mango) GetNeedFetchBaseInfoUrl() show_detail_model.DetailUrls {
	urls := getNeedFetchBaseInfoUrl(m.PlatformId)

	for k, v := range urls {
		urls[k].Url = m.filterUrl(v.Url)
	}

	return urls
}

func (m *mango) filterUrl(u string) string {
	if strings.Contains(u, "/h/") {
		cInfo := &core.CollectorInfo{
			"芒果tv",
			m.Host,
			m.PlatformId,
			"芒果tv介绍页转详情页",
		}

		uArr := strings.Split(u, "/")
		lastSeg := uArr[len(uArr)-1]
		cid := strings.Split(lastSeg, ".")[0]
		cb := fmt.Sprintf("jsonp_%d_%d", time.Now().UnixNano()/1e6, time.Now().Unix()/1e5)
		apiUrl := fmt.Sprintf("https://pcweb.api.mgtv.com/episode/list?_support=10000000&collection_id=%s&size=30&callback=%s",
			cid,
			cb,
		)
		// 爬虫配置
		c := core.NewCollector(cInfo) // 初始化爬虫

		extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
		c.UseProxy()                            // 使用代理，如果代理没有则不使用
		c.DetectCharset = true                  // 非utf-8字符集支持

		c.OnError(func(r *colly.Response, e error) {
			c.Retry(r, e)
		})

		link := ""
		c.OnResponse(func(r *colly.Response) {
			filtered := bytes.Trim(r.Body, cb+"();")
			linkTmp, _ := jsonparser.GetString(filtered, "data", "list", "[0]", "url")
			if linkTmp == "" {
				log.Warn(fmt.Sprintf("未能解析芒果TV的link，源连接：%s, 源数据：%s", u, string(filtered)))
				return
			}
			link = "https://www.mgtv.com/" + strings.TrimLeft(linkTmp, "/")
		})

		c.Visit(apiUrl)
		return link
	} else if strings.Contains(u, "/v/") {
		uParse := strings.Split(u, "/")
		cid := uParse[len(uParse)-3]
		uEnd := uParse[len(uParse)-1]
		uParse = strings.Split(uEnd, ".")
		vid := uParse[0]
		link := fmt.Sprintf("%s/b/%s/%s.html",
			m.Host,
			vid,
			cid,
		)
		return link
	} else {
		return u
	}
}

// 存储播放量
func (m *mango) StorePlayCount(pc int64, ja uint, sid uint64) {
	play_count_model.StorePlayCount(pc, ja, sid, m.PlatformId)
	//fmt.Println(pc)
	//db.GetCrawlerMysqlConn().Create()
}

// 存储评论数量
func (m *mango) StoreCommentCount(cc int64, jobAt uint, showId uint64) {
	comment_count_model.SaveCommentCount(cc, jobAt, showId, m.PlatformId)
}

// 存储艺人信息
func (m *mango) StoreBaseInfoMap(bim map[string]string, showId uint64) {
	storeBaseInfo(bim, showId)
}

// 存储弹幕数
func (m *mango) StoreDanmakuCount(dc int64, ja uint, sid uint64) {
	danmaku_count_model.SaveDanmakuCount(dc, ja, sid, m.PlatformId)
}

// 存储弹幕内容
func (m *mango) StoreDanmakuContent(cts []string, ja uint, sid uint64, ctsId []string) int64 {

	return danmaku_model.SaveDanmaku(cts, ja, sid, m.PlatformId, ctsId)
}
