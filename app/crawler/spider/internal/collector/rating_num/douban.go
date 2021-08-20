package rating_num

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strconv"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func doubanHandle() {

	detailUrls := storage.Douban.GetDetailUrl()


	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go doubanRatingNum(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取评分
func doubanRatingNum(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()

	c := common.Douban.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	c.OnHTML("#interest_sectl > div > div.rating_self.clearfix > strong", func(ele *colly.HTMLElement) {
		rnTemp := ele.Text
		rn := doubanFilterRatingNum(rnTemp)
		findFlag = true
		storage.Douban.StoreRatingNum(rn, JobAt, sid)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到评分", u))
	}
}

func doubanFilterRatingNum(s string) float64 {
	rnFloat, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return rnFloat
}
