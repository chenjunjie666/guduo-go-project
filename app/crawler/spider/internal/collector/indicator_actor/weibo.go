package indicator_actor

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func weiboHandle() {
	urls := storage.Weibo.GetIndicatorActorUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go weiboIndicator(row.Name, row.Id)
	}
	wg.Done()
}

func weiboIndicator(name string, actorId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	u := common.Weibo.FetchIndicatorSearchResUrl(name)

	c := common.Weibo.MobileCollector(ModName)

	wid := common.Weibo.ParseWid(u)

	if wid == "" {
		log.Warn(fmt.Sprintf("微博指数获取失败，分析wid失败，wid为空，源数据为：%s", name))
		return
	}

	dataGroup := "1hour"
	apiUrl := fmt.Sprintf("%s/index/ajax/newindex/getchartdata?wid=%s&dateGroup=%s",
		common.Weibo.ApiHosts.DataHost,
		wid,
		dataGroup,
	)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false
	c.OnResponse(func(r *colly.Response) {
		j := r.Body

		idx := 0
		idx2 := 0
		_, _ = jsonparser.ArrayEach(j, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if idx != 0 {
				return
			}
			idx++
			_, _ = jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				if idx2 == 12 {
					in, e := jsonparser.ParseInt(value)
					if e != nil {
						log.Warn(fmt.Sprintf("解析微博指数失败：%s，源数据：%s", e, string(j)))
						return
					}

					findFlag = true
					storage.Weibo.StoreIndicatorActor(in, JobAt, actorId)
				}
				idx2++
			}, "trend", "s")
		}, "data")

	})

	// 设置 referer 参数 防止CSRF
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("referer", u)
	})

	_ = c.Visit(apiUrl)

	if !findFlag {
		log.Warn(fmt.Sprintf("获取微博指数失败：%s", u))
	}
}
