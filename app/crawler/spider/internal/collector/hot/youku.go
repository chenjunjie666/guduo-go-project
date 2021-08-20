package hot

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strconv"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func youkuHandle() {
	detailUrls := storage.Youku.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go youkuHot(row.Url, row.ShowId)
	}
	wg.Done()
}

// 抓取优酷热度趋势
func youkuHot(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Youku.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false
	// 找到html中的热度class
	c.OnHTML(".video-heat-text", func(ele *colly.HTMLElement) {
		hotStr := ele.Text                         // 获取html内容
		hot, e := strconv.ParseInt(hotStr, 10, 64) // string转int64
		if e != nil {
			log.Warn(fmt.Sprintf("优酷热度执行数据转换失败，原因：%s，源数据：%s", e, hotStr))
		}

		findFlag = true
		storage.Youku.StoreHot(hot, JobAt, showId) // 保存热度
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到热度趋势", u))
	}
}
