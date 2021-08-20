package indicator_actor

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func baiduHandle() {
	urls := storage.Baidu.GetIndicatorActorUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go baiduIndicator(row.IndicatorUrl, row.ActorId)
	}
	wg.Done()
}

// 爬取百度指数
func baiduIndicator(u string, actorId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	apiHost := common.Baidu.ApiHosts.IndexHost
	words := common.Baidu.ParseWords(u)
	// 经过观察以下参数是固定值
	wordType := 1
	days := 2
	area := 0

	// 格式化获取百度指数的链接
	u = fmt.Sprintf("%s/api/SearchApi/index?word=[[{\"name\":\"%s\",\"wordType\":%d}]]&days=%d&area=%d",
		apiHost, // 网站host
		words,
		wordType,
		days, // 时间跨度参数
		area,
	)

	c := common.Baidu.Collector(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false

	c.OnResponse(func(r *colly.Response) {
		msg, _ := jsonparser.GetString(r.Body, "message")
		if msg == "bad request" {
			c.Retry(r, nil)
			return
		}

		findFlag = true
		ctxStr := string(r.Body)
		ctxByte := []byte(ctxStr)
		uniqueId, _ := jsonparser.GetString(ctxByte, "data", "uniqid")
		encryptData, _ := jsonparser.GetString(ctxByte, "data", "userIndexes", "[0]", "all", "data")
		decryptData := ""
		lastDayDataInt := int64(0)
		ptbk := common.Baidu.ParsePtbk(uniqueId)
		if ptbk != "" {
			decryptData = common.Baidu.DecryptData(ptbk, encryptData)
			lastDayDataInt = common.Baidu.ParseLastDayData(decryptData)
		}

		log.Info("++++++++++++++++++++++++百度演员指数爬取完毕，actor id:", actorId, " 抓取完毕，值为：" ,lastDayDataInt , ",准备开始存储++++++++++++++++++++++++")
		storage.Baidu.StoreIndicatorActor(lastDayDataInt, JobAt, actorId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到百度指数", u))
	}
}
