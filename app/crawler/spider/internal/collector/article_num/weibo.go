package article_num

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func weiboHandle() {
	detailUrls := storage.Weibo.GetFetchArticleUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go weiboArticleNum(row.Url, row.ID)
	}

	wg.Done()
}

// 从页面解析并获取获取微博文章数
func weiboArticleNum(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Weibo.CollectorWithLogin(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", u)
	})

	findFlag := false
	c.OnHTML("#pl_feedlist_index > div.m-error", func(ele *colly.HTMLElement) {
		findFlag = true
		anTemp := ele.Text
		an := filterArticleNum(anTemp)
		storage.Wechat.StoreArticleNum(an, showId, JobAt)
	})

	//c.OnResponse(func(r *colly.Response) {
	//	fmt.Println(string(r.Body))
	//})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到文章数", u))
	}

}
