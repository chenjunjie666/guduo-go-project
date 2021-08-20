package danmaku

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core"
	"strings"
	"testing"
	"time"
)

func TestMangoDanmakuContent(t *testing.T) {
	core.Init()
	ch.PushJob()
	wg.Add(1)
	//url := storage.Mango.GetDetailUrl()[0]
	//mangoDanmakuContent("https://www.mgtv.com/b/360586/10984639.html?fpa=127&fpos=2&lastp=ch_cartoon&cpid=8", url.ShowId)
	url := "https://www.mgtv.com/b/368455/11861742.html"
	mangoDanmakuContent(filterMangoUrl(url), 22298)
}



func filterMangoUrl(u string) string {
	if strings.Contains(u, "/h/") {
		cInfo := &core.CollectorInfo{
			"芒果tv",
			common.Mango.Host,
			common.Mango.PlatformId,
			"芒果tv详情页",
		}

		uArr := strings.Split(u, "/")
		lastSeg := uArr[len(uArr) - 1]
		cid := strings.Split(lastSeg, ".")[0]
		cb := fmt.Sprintf("jsonp_%d_%d", time.Now().UnixNano() / 1e6, time.Now().Unix() / 1e5)
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
			filtered := bytes.Trim(r.Body, cb + "();")
			linkTmp, _ := jsonparser.GetString(filtered, "data", "list", "[0]", "url")
			if linkTmp == "" {
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
			common.Mango.Host,
			vid,
			cid,
		)
		return link
	} else {
		return u
	}
}

