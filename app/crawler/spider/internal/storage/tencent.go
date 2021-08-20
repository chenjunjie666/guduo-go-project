package storage

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/internal/constant"
	comment_count_model "guduo/app/internal/model_scrawler/comment_count_model"
	"guduo/app/internal/model_scrawler/danmaku_count_model"
	danmaku_model "guduo/app/internal/model_scrawler/danmaku_model"
	play_count_model "guduo/app/internal/model_scrawler/play_count_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
	show_model "guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"strings"
	"sync"
)

var Tencent = &tencent{
	PlatformId: constant.PlatformIdTencent,
	Host:       "https://v.qq.com",
}

type tencent struct {
	PlatformId uint64
	Host       string
}

func (t tencent) GetNeedFetchIntroUrl() show_detail_model.DetailUrls {
	f := show_model.GetActiveShows()
	sids := make([]model.PrimaryKey, len(f))
	for k, v := range f {
		sids[k] = v
	}

	urls := show_detail_model.GetDetailUrl(t.PlatformId, sids)

	for k, v := range urls {
		urls[k].Url = t.filterUrl(v.Url)
	}

	return urls
}

// 获取腾讯视频剧集的详情页链接
func (t *tencent) GetDetailUrl() show_detail_model.DetailUrls {
	//showIds := show_model.GetActiveShows()
	//urls := show_detail_model.GetDetailUrl(t.PlatformId, showIds)
	urls := show_detail_model.GetDetailUrlNew(t.PlatformId)

	log.Info("腾讯视频，一共", len(urls), "条视频，需要爬取")
	wg2 := &sync.WaitGroup{}
	for k, v := range urls {
		oUrl := v.Url
		wg2.Add(1)
		go func(k int, sid uint64) {
			urls[k].Url = t.filterUrl(oUrl)
			if urls[k].Url != oUrl {
				show_detail_model.SaveTrueUrl(urls[k].Url, sid, constant.PlatformIdTencent)
			}
			wg2.Done()
		}(k, v.ShowId)
	}
	wg2.Wait()

	return urls
	//urls := []string{
	//	"https://v.qq.com/x/cover/g9fv8x19fzl9id9/a0026t8wfcc.html",
	//	"https://v.qq.com/x/cover/rjae621myqca41h.html",
	//}
	//return urls
}

func (t tencent) GetNeedFetchBaseInfoUrl() show_detail_model.DetailUrls {
	urls := getNeedFetchBaseInfoUrl(t.PlatformId)

	for k, v := range urls {
		urls[k].Url = t.filterUrl(v.Url)
	}

	return urls
}

func (t tencent) filterUrl(u string) string {
	if !strings.Contains(u, "detail") {
		return u
	}
	cInfo := &core.CollectorInfo{
		"腾讯视频",
		t.Host,
		t.PlatformId,
		"腾讯视频详情页",
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
	c.OnHTML(".mod_episode .item:first-child a", func(ele *colly.HTMLElement) {
		link = ele.Attr("href")
	})
	c.Visit(u)

	if !strings.Contains(link, "qq.com/x/cove") {
		link = ""
	}

	if link == "" {
		log.Info("腾讯连接", u, "未找到播放页连接")
		return ""
	}

	if !strings.Contains(link, "http") {
		link = "https:" + link
	}

	return link
}

func (t tencent) StoreIntro(intro string, showId uint64) {
	storeIntro(intro, showId)
}

// 存储播放量
func (t *tencent) StorePlayCount(pc int64, ja uint, sid uint64) {
	play_count_model.StorePlayCount(pc, ja, sid, t.PlatformId)

	//pcM := &play_count_model.Def{
	//}

	//fmt.Println(pc)
	//db.GetCrawlerMysqlConn().Create()
}

// 存储艺人信息
func (t *tencent) StoreBaseInfoMap(bim map[string]string, showId uint64) {
	storeBaseInfo(bim, showId)
}

// 存储评论数
func (t *tencent) StoreCommentCount(cc int64, jobAt uint, showId uint64) {
	comment_count_model.SaveCommentCount(cc, jobAt, showId, t.PlatformId)
}

// 存储弹幕数
func (t *tencent) StoreDanmakuCount(dc int64, ja uint, sid uint64) {
	danmaku_count_model.SaveDanmakuCount(dc, ja, sid, t.PlatformId)
}

// 存储弹幕内容
func (t *tencent) StoreDanmakuContent(dcMap []string, dcIdArr[]string, ja uint, sid uint64) int64 {
	return danmaku_model.SaveDanmaku(dcMap, ja, sid, t.PlatformId, dcIdArr)
}
