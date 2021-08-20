package length

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func doubanHandle() {

	detailUrls := storage.Douban.GetDetailUrl()

	wg.Add(len(detailUrls))
	for _, row := range detailUrls {
		go doubanLength(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取片长
func doubanLength(u string, sid uint64) {
	defer wg.Done()

	c := common.Douban.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	c.OnHTML("#info", func(ele *colly.HTMLElement) {
		findFlag = true
		lenTemp := ele.Text

		lenStr := filterLength(lenTemp)
		storage.Douban.StoreLength(lenStr, sid)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到片长", u))
	}
}

func filterLength(s string) string {
	reg := regexp.MustCompile(`\d+分钟`)
	result := reg.FindAllStringSubmatch(s, -1)
	str := ""
	count := 1
	for _, text := range result {
		str += text[0]
		if count < len(result) {
			count += 1
			str += "/"
		}
	}
	return str
}
