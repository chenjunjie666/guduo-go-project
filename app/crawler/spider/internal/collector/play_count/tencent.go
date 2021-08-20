// 腾讯剧集播放量指标采集
package play_count

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

// 腾讯播放量采集指标入口
func tencentHandle() {
	detailUrls := storage.Tencent.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go tencentPlayCount(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取播放量
func tencentPlayCount(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Tencent.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	c.OnResponse(func(r *colly.Response) {
		reg := regexp.MustCompile(`"view_all_count":\d+`)
		pcString := reg.FindString(string(r.Body))
		pcString = strings.Trim(pcString, `"view_all_count":`)
		pc, e := strconv.ParseInt(pcString, 10, 64)

		if e != nil {
			log.Warning(fmt.Sprintf("找到播放量但是解析失败：%s, 连接为:%s", e, u))
			return
		}

		findFlag = true
		storage.Tencent.StorePlayCount(pc, JobAt, sid)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到播放量", u))
	}
}

// 对播放量进行转换
//func tencentFilterPlayCount(s string) int64 {
//	reg := regexp.MustCompile(`([亿万])`)
//
//	r := reg.Find([]byte(s))
//
//	multi := float64(1)
//	if len(r) > 0 {
//		switch string(r) {
//		case "亿":
//			multi = 1e9
//		case "万":
//			multi = 1e5
//		}
//		s = strings.TrimRight(s, string(r))
//	}
//
//	pcFloat, err := strconv.ParseFloat(s, 64)
//	pcFloat = pcFloat * multi
//	pcInt := int64(pcFloat)
//
//	if err != nil {
//		return 0
//	}
//
//	return pcInt
//}
