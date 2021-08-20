package indicator

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	guduoJson "guduo/pkg/json"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func sogouHandle() {
	urls := storage.Sogou.GetIndicatorUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go sogouIndicator(row.Url, row.ShowId)
	}
	wg.Done()
}

// 爬取搜狗指数
func sogouIndicator(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	apiHost := common.Sogou.ApiHosts.IndexHost
	words := common.Sogou.ParseWords(u)
	startDate := common.Sogou.ParseStartDate()
	endDate := common.Sogou.ParseEndDate()
	// 经过观察以下参数是固定值
	dataType := "SEARCH_ALL"
	queryType := "INPUT"

	// 格式化获取搜狗指数的链接
	u = fmt.Sprintf("%s/getDateData?kwdNamesStr=%s&startDate=%s&endDate=%s&dataType=%s&queryType=%s",
		apiHost, // 网站host
		words,
		startDate,
		endDate,
		dataType,
		queryType,
	)

	c := common.Sogou.Collector(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false

	c.OnResponse(func(r *colly.Response) {

		findFlag = true
		ctxStr := string(r.Body)
		ctxByte := []byte(ctxStr)
		jsonArrayLength, _ := guduoJson.JsonParserGetArrayLength(ctxByte, "data", "pvList", "[0]")
		targetIndex := fmt.Sprintf("[%d]", jsonArrayLength-2)
		lastDayDataInt, _ := jsonparser.GetInt(ctxByte, "data", "pvList", "[0]", targetIndex, "pv")

		storage.Sogou.StoreIndicator(lastDayDataInt, JobAt, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到搜狗指数", u))
	}
}
