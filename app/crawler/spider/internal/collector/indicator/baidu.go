package indicator

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func baiduHandle() {
	urls := storage.Baidu.GetIndicatorUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go baiduGenderIndicator(row.IndicatorUrl, row.ShowId)

		wg.Add(1)
		ch.PushJob()
		go baiduAgeIndicator(row.GenAgeIndicatorUrl, row.ShowId)

		wg.Add(1)
		ch.PushJob()
		go baiduIndicator(row.GenAgeIndicatorUrl, row.ShowId)
	}
	wg.Done()
}

// 爬取百度指数-性别分布
func baiduGenderIndicator(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()

	apiHost := common.Baidu.ApiHosts.IndexHost
	words := common.Baidu.ParseWords(u)

	// 格式化获取性别分布的链接
	u = fmt.Sprintf("%s/api/SocialApi/baseAttributes?wordlist[]=%s",
		apiHost, // 网站host
		words,   // keywords 参数
	)

	c := common.Baidu.CollectorWithoutLogin(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false
	var genderRateMap map[string]float64
	genderRateMap = make(map[string]float64)

	c.OnResponse(func(r *colly.Response) {
		ctxStr := string(r.Body)
		ctxByte := []byte(ctxStr)

		_, _ = jsonparser.ArrayEach(ctxByte, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			desc, _ := jsonparser.GetString(value, "desc")
			rate, _ := jsonparser.GetFloat(value, "rate")
			genderRateMap[desc] = rate
		}, "data", "result", "[0]", "gender")

		res := make(map[string]float64)
		res["male"] = genderRateMap["男"]
		res["female"] = genderRateMap["女"]
		findFlag = true
		log.Info("++++++++++++++++++++++++百度性别分布爬取完毕，showid:", showId, " 抓取完毕，准备开始存储++++++++++++++++++++++++")
		storage.Baidu.StoreGenderRateMap(res, JobAt, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("showid:%d, 链接:%s，未找到性别分布", showId, u))
	}

}

// 爬取百度指数-年龄分布
func baiduAgeIndicator(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	apiHost := common.Baidu.ApiHosts.IndexHost
	words := common.Baidu.ParseWords(u)

	// 格式化获取年龄分布的链接
	u = fmt.Sprintf("%s/api/SocialApi/baseAttributes?wordlist[]=%s",
		apiHost, // 网站host
		words,   // keywords 参数
	)

	c := common.Baidu.CollectorWithoutLogin(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false
	var ageRateMap map[string]float64
	ageRateMap = make(map[string]float64)

	c.OnResponse(func(r *colly.Response) {
		ctxStr := string(r.Body)
		ctxByte := []byte(ctxStr)

		_, _ = jsonparser.ArrayEach(ctxByte, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			desc, _ := jsonparser.GetString(value, "desc")
			rate, _ := jsonparser.GetFloat(value, "rate")
			ageRateMap[desc] = rate
		}, "data", "result", "[0]", "age")
		if len(ageRateMap) == 0 {
			return
		}
		findFlag = true
		log.Info("++++++++++++++++++++++++百度年龄分布爬取完毕，showid:", showId, " 抓取完毕，准备开始存储++++++++++++++++++++++++")
		storage.Baidu.StoreAgeRateMap(ageRateMap, JobAt, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("showid：%d, 链接:%s，未找到年龄分布", showId, u))
	}
}

// 爬取百度指数
func baiduIndicator(u string, showId uint64) {
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

		if lastDayDataInt == 0 {
			log.Warn("百度指数未获取到，链接为", u, "源数据为", ctxStr)
		}
		fmt.Println(lastDayDataInt)
		findFlag = true
		log.Info("++++++++++++++++++++++++百度指数爬取完毕，showid:", showId, " 抓取完毕，准备开始存储++++++++++++++++++++++++")
		storage.Baidu.StoreIndicator(lastDayDataInt, JobAt, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("showid：%d, 链接:%s，未找到百度指数", showId, u))
	}
}
