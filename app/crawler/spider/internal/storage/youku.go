package storage

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/internal/constant"
	comment_count_model "guduo/app/internal/model_scrawler/comment_count_model"
	danmaku_count_model "guduo/app/internal/model_scrawler/danmaku_count_model"
	danmaku_model "guduo/app/internal/model_scrawler/danmaku_model"
	hot_model "guduo/app/internal/model_scrawler/hot_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
	"strings"
	"sync"
)

var Youku = &youku{
	PlatformId: constant.PlatformIdYouku,
	Host:       "https://v.youku.com",
}

type youku struct {
	PlatformId uint64
	Host       string
}

// 获取优酷的详情页链接
func (y *youku) GetDetailUrl() show_detail_model.DetailUrls {
	//showIds := show_model.GetActiveShows()
	//urls := show_detail_model.GetDetailUrl(y.PlatformId, showIds)
	urls := show_detail_model.GetDetailUrlNew(y.PlatformId)

	log.Info("优酷视频，一共", len(urls), "条视频，需要爬取")

	wg2 := &sync.WaitGroup{}
	for k, v := range urls {
		oUrl := v.Url
		wg2.Add(1)
		go func(k int, sid uint64) {
			urls[k].Url = y.filterLinks(oUrl)
			if urls[k].Url != oUrl {
				show_detail_model.SaveTrueUrl(urls[k].Url, sid, constant.PlatformIdYouku)
			}
			wg2.Done()
		}(k, v.ShowId)
	}
	wg2.Wait()
	return urls

	//urls := []string{
	//	"https://v.youku.com/v_show/id_XMjY3MTQ2MDE0OA==.html",
	//	"https://v.youku.com/v_show/id_XNDk2MzA5Nzc2MA==.html",
	//}
	//return urls
}

func (y *youku) GetNeedFetchBaseInfoUrl() show_detail_model.DetailUrls {
	urls := getNeedFetchBaseInfoUrl(y.PlatformId)
	for k, v := range urls {
		urls[k].Url = y.filterLinks(v.Url)
	}
	return urls
}

func (y *youku) filterLinks(u string) string {
	if !strings.Contains(u, "list") {
		return u
	}
	cInfo := &core.CollectorInfo{
		"优酷",
		y.Host,
		y.PlatformId,
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
	})

	c.Visit(u)
	if !strings.Contains(link, "youku.com/v_show") {
		link = ""
	}

	if link == "" {
		log.Info("优酷连接", u, "未找到播放页连接")
		return ""
	}

	if !strings.Contains(link, "http") {
		link = "https:" + link
	}
	return link
}

// 存储优酷热度
func (y *youku) StoreHot(hot int64, jobAt uint, showId uint64) {
	hot_model.SaveHotCount(hot, jobAt, showId, y.PlatformId)
}

// 存储艺人信息
func (y *youku) StoreBaseInfoMap(bim map[string]string, showId uint64) {
	storeBaseInfo(bim, showId)
}

func (y *youku) StoreCommentCount(cc int64, ja uint, sid uint64) {
	comment_count_model.SaveCommentCount(cc, ja, sid, y.PlatformId)
}

func (y *youku) StoreDanmankuContent(dmk []string, ja uint, sid uint64) {
	danmaku_model.SaveDanmaku(dmk, ja, sid, y.PlatformId)
}

func (y *youku) StoreDanmankuCount(dc int64, ja uint, sid uint64) {
	danmaku_count_model.SaveDanmakuCount(dc, ja, sid, y.PlatformId)
}
