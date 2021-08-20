package article_num_actor

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func wechatHandle() {
	detailUrls := storage.Wechat.GetArticleNumActorUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go wechatArticleNum(row.Url, row.ID)
	}

	wg.Done()
}

// 从页面解析并获取获取微信文章数
func wechatArticleNum(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Wechat.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Refer", "https://weixin.sogou.com/")
	})

	c.OnHTML("#pagebar_container > div.mun", func(ele *colly.HTMLElement) {
		findFlag = true
		anTemp := ele.Text
		an := filterArticleNum(anTemp)
		storage.Wechat.StoreArticleNumActor(an, showId, JobAt)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到文章数", u))
	}
}

func filterArticleNum(s string) int64 {
	reg := regexp.MustCompile(`\d+`)
	result := reg.FindAllString(s, -1)
	str := ""
	for _, text := range result {
		str += text
	}
	anInt, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return 0
	}

	return anInt
}
