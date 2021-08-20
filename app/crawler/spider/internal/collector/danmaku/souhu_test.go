package danmaku

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core"
	"regexp"
	"strings"
	"testing"
)

func TestSouhuDanmakuContent(t *testing.T) {
	core.Init()

	wg.Add(1)
	//url := storage.Souhu.GetDetailUrl()[0]
	//souhuDanmakuContent("https://film.sohu.com/album/9702422.html", 0)

	u := "https://film.sohu.com/album/9440003.html"
	souhuDanmakuContent(filerSohuUrl(u), 22212)
}

func TestBingfa(t *testing.T) {
	common.Souhu.ParseVids("https://tv.sohu.com/v/MjAyMTA0MTUvbjYwMDk5Njk5MC5zaHRtbA==.html")
}





func filerSohuUrl(u string) string {

	if strings.Contains(u, "com/s") {
		cInfo := &core.CollectorInfo{
			"搜狐TV",
			common.Souhu.Host,
			common.Souhu.PlatformId,
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
			common.Souhu.Host,
			common.Souhu.PlatformId,
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
			common.Souhu.Host,
			common.Souhu.PlatformId,
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
