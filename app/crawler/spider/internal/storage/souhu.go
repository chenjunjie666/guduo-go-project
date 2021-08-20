package storage

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_scrawler/comment_count_model"
	"guduo/app/internal/model_scrawler/danmaku_count_model"
	"guduo/app/internal/model_scrawler/danmaku_model"
	"guduo/app/internal/model_scrawler/show_detail_model"
	"regexp"
	"strings"
	"sync"
)

var Souhu = &souhu{
	PlatformId: constant.PlatformIdSouhu,
	Host:       "https://tv.sohu.com/",
}

type souhu struct {
	PlatformId uint64
	Host       string
}

// 获取搜狐视频剧集的详情页链接
func (s *souhu) GetDetailUrl() show_detail_model.DetailUrls {
	//showIds := show_model.GetActiveShows()
	//urls := show_detail_model.GetDetailUrl(s.PlatformId, showIds)
	urls := show_detail_model.GetDetailUrlNew(s.PlatformId)

	log.Info("搜狐视频，一共", len(urls), "条视频，需要爬取")

	wg2 := &sync.WaitGroup{}
	for k, v := range urls {
		oUrl := v.Url
		wg2.Add(1)
		go func(k int, sid uint64) {
			urls[k].Url = s.filterUrl(oUrl)
			if urls[k].Url != oUrl {
				show_detail_model.SaveTrueUrl(urls[k].Url, sid, constant.PlatformIdSouhu)
			}
			wg2.Done()
		}(k, v.ShowId)
	}
	wg2.Wait()

	return urls
	//urls := []string{
	//	"https://tv.sohu.com/v/MjAxNzA2MDUvbjQ5NTc5NjM4Ny5zaHRtbA==.html",
	//	"https://tv.sohu.com/v/MjAyMTAzMjIvbjYwMDk5MDQ4MS5zaHRtbA==.html",
	//}
	//return urls
}
func (s souhu) GetNeedFetchBaseInfoUrl() show_detail_model.DetailUrls {
	urls := getNeedFetchBaseInfoUrl(s.PlatformId)

	for k, v := range urls {
		urls[k].Url = s.filterUrl(v.Url)
	}

	return urls
}

func (s souhu) filterUrl(u string) string {
	if strings.Contains(u, "com/s") {
		cInfo := &core.CollectorInfo{
			"搜狐TV",
			s.Host,
			s.PlatformId,
			"搜狐TV搜索视频播放页",
		}
		// 爬虫配置
		c := core.NewCollector(cInfo)           // 初始化爬虫
		extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent

		link := ""
		playlistId := ""
		c.OnResponse(func(r *colly.Response) {
			html := string(r.Body)
			//fmt.Println(html)

			reg := regexp.MustCompile(`playlistId = "\d+"`)
			reg2 := regexp.MustCompile(`PLAYLIST_ID = "\d+"`)

			playlistId1 := reg.FindString(html)
			playlistId2 := reg2.FindString(html)

			if playlistId1 != "" {
				playlistId = strings.Trim(playlistId1, `playlistId=" `)
			} else if playlistId2 != "" {
				playlistId = strings.Trim(playlistId2, `PLAYLIST_ID=" `)
			}

			if playlistId == "" {
				return
			}
		})

		c.Visit(u)
		if playlistId == "" {
			log.Warning("连接", u, "未找到播放页")
			return ""
		}

		// 爬虫配置
		c = core.NewCollector(cInfo)            // 初始化爬虫
		extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent

		c.OnResponse(func(r *colly.Response) {
			html := string(r.Body)

			reg := regexp.MustCompile(`tvUrl:".*?"`)
			tvUrl := reg.FindString(html)

			link = strings.Trim(tvUrl, `tvUrl:"`)

			if !strings.Contains(link, "sohu.com") {
				link = ""
			}
		})
		mUrl := fmt.Sprintf("https://m.tv.sohu.com/album/s%s.shtml", playlistId)
		c.Visit(mUrl)
		if link == "" {
			log.Warning("连接", u, "未找到播放页")
			return ""
		}
		if !strings.Contains(link, "http") {
			link = "https:" + link
		}

		if !strings.Contains(link, "html") {
			link = link + "l"
		}

		return link
	} else if strings.Contains(u, "so.tv") {
		cInfo := &core.CollectorInfo{
			"搜狐TV",
			s.Host,
			s.PlatformId,
			"搜狐TV搜索视频播放页",
		}
		// 爬虫配置
		c := core.NewCollector(cInfo)           // 初始化爬虫
		extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent

		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Cookie", "SUV=1905271314163369")
		})

		link := ""

		c.OnHTML(".pic:first-child a:first-child", func(ele *colly.HTMLElement) {
			link = ele.Attr("href")
			if !strings.Contains(link, "sohu.com") {
				link = ""
			}
		})

		c.Visit(u)
		if link == "" {
			log.Warning("连接", u, "未找到播放页")
			return ""
		}
		if !strings.Contains(link, "http") {
			link = "https:" + link
		}

		if !strings.Contains(link, "html") {
			link = link + "l"
		}

		return link
	} else if strings.Contains(u, "/item/") {
		cInfo := &core.CollectorInfo{
			"搜狐TV",
			s.Host,
			s.PlatformId,
			"搜狐TV搜索视频播放页",
		}
		// 爬虫配置
		c := core.NewCollector(cInfo)           // 初始化爬虫
		extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent

		link := ""
		c.OnHTML(".colL div a:first-child", func(ele *colly.HTMLElement) {
			link = ele.Attr("href")
			if !strings.Contains(link, "sohu.com") {
				link = ""
			}
		})

		c.Visit(u)
		if link == "" {
			log.Warning("连接", u, "未找到播放页")
			return ""
		}
		if !strings.Contains(link, "http") {
			link = "https:" + link
		}

		if !strings.Contains(link, "html") {
			link = link + "l"
		}

		return link
	}

	if strings.Contains(u, "http") && !strings.Contains(u, "https") {
		u = strings.Replace(u, "http://", "https://", 1)
	}

	return u
}

// 存储播放量
//func (s *souhu) StorePlayCount(pc int64) {
//fmt.Println(pc)
//db.GetCrawlerMysqlConn().Create()
//}

// 存储评论数量
func (s *souhu) StoreCommentCount(cc int64, jobAt uint, showId uint64) {
	comment_count_model.SaveCommentCount(cc, jobAt, showId, s.PlatformId)
}

// 存储艺人信息
func (s *souhu) StoreBaseInfoMap(bim map[string]string, showId uint64) {
	storeBaseInfo(bim, showId)
}

// TODO
// 存储弹幕数
func (s *souhu) StoreDanmakuCount(dc int64, ja uint, sid uint64) {
	danmaku_count_model.SaveDanmakuCount(dc, ja, sid, s.PlatformId)
}

// 存储弹幕内容
func (s *souhu) StoreDanmakuContent(dcMap []string, ja uint, sid uint64) {

	danmaku_model.SaveDanmaku(dcMap, ja, sid, s.PlatformId)
}
