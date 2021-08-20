package news_num

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func baiduHandle() {
	detailUrls := storage.Baidu.GetNewsUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go baiduNewsNum(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取百度新闻资讯数
func baiduNewsNum(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()

	c := common.Baidu.CollectorWithoutLogin(ModName)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "www.baidu.com")
	})

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	c.OnHTML(".nums", func(ele *colly.HTMLElement) {
		findFlag = true
		nnTemp := ele.Text
		nn := filterNewsNum(nnTemp)
		storage.Baidu.StoreNewsNum(nn, JobAt, sid)
	})

	html := ""
	c.OnResponse(func(r *colly.Response) {
		html = string(r.Body)
	})

	_ = c.Visit(u)

	if !findFlag {
		fmt.Println(html)
		log.Warn(fmt.Sprintf("链接:%s，未找到新闻资讯数", u))
	}
}

func filterNewsNum(s string) int64 {
	reg := regexp.MustCompile(`\d+`)
	result := reg.FindAllString(s, -1)
	str := ""
	for _, text := range result {
		str += text
	}
	nnInt, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return 0
	}

	return nnInt
}
