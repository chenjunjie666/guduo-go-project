package release_time

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func baiduHandle() {

	detailUrls := storage.Baidu.GetBaikeUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go baiduReleaseTime(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取上映日期
func baiduReleaseTime(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()

	c := common.Baidu.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	c.OnHTML("div.basic-info.cmn-clearfix", func(ele *colly.HTMLElement) {
		rtTemp := ele.Text
		rt := filterReleaseTime(rtTemp)
		rts := formatReleaseTime(rt)
		if rts <= 1000000000 {
			return
		}
		findFlag = true
		storage.Baidu.StoreReleaseTimeStamp(rts, sid)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到上映日期", u))
	}
}

func filterReleaseTime(s string) string {
	reg := regexp.MustCompile(`(上映时间|首播时间)(\s|\n|\r|\r\n|\n\r)*\d+年\d+月\d+日?`)
	str := reg.FindString(s)
	return str
}

func formatReleaseTime(datetimeStr string) uint {
	if datetimeStr == "" {
		return 0
	}
	//日期转化为时间戳
	//Go语言中格式化时间模板不是常见的Y-m-d H:M:S而是使用Go语言的诞生时间 2006-01-02 15:04:05 -0700 MST
	reg := regexp.MustCompile(`\d+`)
	result := reg.FindAllString(datetimeStr, -1)
	datetimeStr = ""
	count := 1
	for _, text := range result {
		if len(text) < 2 {
			text = "0" + text
		}
		datetimeStr += text
		if count < len(result) {
			count += 1
			datetimeStr += "-"
		} else {
			datetimeStr += " 00:00:00"
		}
	}

	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, datetimeStr, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64
	return uint(timestamp)
}
