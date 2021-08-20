package short_comment

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func doubanHandle() {
	urls := storage.Douban.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go doubanShortComment(row.Url, row.ShowId)
	}

	wg.Done()
}

// 爬取豆瓣短评数主逻辑
func doubanShortComment(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Douban.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	isFind := false
	// 匹配短评所在html的css选择器
	c.OnHTML("#comments-section h2 a", func(ele *colly.HTMLElement) {
		scCountTmpStr := ele.Text

		// 将 "全部 xxxx 条" 通过正则筛选出其中的数字部分
		reg := regexp.MustCompile(`\d+`)
		scCountStr := string(reg.Find([]byte(scCountTmpStr)))

		// 将筛选出的数字部分转int
		scCount, e := strconv.ParseInt(scCountStr, 0, 64)
		if e != nil {
			log.Warn(fmt.Sprintf("获取豆瓣短评失败，数字转换失败：%s，源数据:%s", e, scCountStr))
			return
		}

		isFind = true
		storage.Douban.StoreShortCommentCount(scCount, JobAt, sid)
	})

	_ = c.Visit(u)

	if !isFind {
		log.Warn(fmt.Sprintf("没有找到豆瓣短评：%s", u))
	}
}
