package common

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2/extensions"
)

var Sogou = &sogou{
	PlatformId: storage.Sogou.PlatformId,
	Host:       storage.Sogou.Host,
	ApiHosts: struct {
		IndexHost string
	}{
		IndexHost: "http://index.sogou.com",
	},
}

type sogou struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		IndexHost string
	}
}

// 搜狗采集器初始化
func (s sogou) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"搜狗",
		s.Host,
		s.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持

	return c
}

// 解析搜狗指标的搜索开始时间
func (s sogou) ParseStartDate() string {
	now := time.Now()
	year := now.Year()   //年
	month := now.Month() //月
	day := now.Day()     //日
	startDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local).AddDate(0, 0, -7)
	startDateStr := startDate.Format("20060102")

	return startDateStr
}

// 解析搜狗指标的搜索开始时间
func (s sogou) ParseEndDate() string {
	now := time.Now()
	year := now.Year()   //年
	month := now.Month() //月
	day := now.Day()     //日
	endDate := time.Date(year, month, day, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	endDateStr := endDate.Format("20060102")

	return endDateStr
}

// 解析搜狗指标的搜索关键词
func (s sogou) ParseWords(u string) string {
	words := ""
	reg := regexp.MustCompile(`(\?|&)kwdNamesStr=.*?(&|$)`)
	words = reg.FindString(u)
	words = strings.Trim(words, `?kwdNamesStr=&`)

	return words
}

// 根据返回结果集，返回最近一天搜狗指标数
func (s sogou) ParseLastDayData(data string) int64 {
	dataParse := strings.Split(data, ",")
	lastDayData := dataParse[len(dataParse)-1]
	lastDayDataInt, _ := strconv.ParseInt(lastDayData, 10, 64)

	return lastDayDataInt
}
