package introduction

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func tencentHandle() {
	urls := storage.Tencent.GetNeedFetchIntroUrl()

	wg.Add(len(urls))
	for _, row := range urls {
		go tencentIntroduction(row.Url, row.ShowId)
	}
	wg.Done()
}

// 爬取腾讯剧情简介主逻辑
func tencentIntroduction(u string, showId uint64) {
	defer wg.Done()
	c := common.Tencent.Collector(ModName)

	findFlag := false

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	// 匹配 css 规则获取剧情简介内容
	c.OnHTML(".video_summary p", func(ele *colly.HTMLElement) {
		intro := ele.Text
		findFlag = true
		// 存储获取到的简介
		storage.Tencent.StoreIntro(intro, showId)
	})

	_ = c.Visit(u)

	// 如果没有找到，记录错误日志
	if !findFlag {
		log.Warn(fmt.Sprintf("获取剧情简介失败，链接：%s", u))
	}
}
