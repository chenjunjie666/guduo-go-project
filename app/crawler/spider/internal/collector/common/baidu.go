package common

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/crawler/spider/internal/util"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

var Baidu = &baidu{
	PlatformId: storage.Baidu.PlatformId,
	Host:       storage.Baidu.Host,
	ApiHosts: struct {
		IndexHost string
		TiebaHost string
		BaikeHost string
		NewsHost  string
	}{
		IndexHost: "https://index.baidu.com",
		TiebaHost: "https://tieba.baidu.com",
		BaikeHost: "https://baike.baidu.com",
		NewsHost:  "https://news.baidu.com",
	},
}

type baidu struct {
	PlatformId uint64
	Host       string
	ApiHosts   struct {
		IndexHost string
		TiebaHost string
		BaikeHost string
		NewsHost  string
	}
}

// 百度采集器初始化
func (b baidu) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"百度",
		b.Host,
		b.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	// 设置 cookie
	cookies := b.getCookies()

	// 设置 cookie
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", util.BuildCookie(cookies))
	})

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}


// 百度采集器初始化
func (b baidu) CollectorWithoutLogin(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"百度",
		b.Host,
		b.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}

func (b baidu) getCookies() []http.Cookie {
	cookies := []http.Cookie{
		{
			Name:  "BDUSS",
			Value: "Ug5MTdxZFRUSE1UVEZkT2pjflVtb3NIYkZjQlBid2dhRlNYMllFa1Bra0wzdDFnSVFBQUFBJCQAAAAAAAAAAAEAAACU7YMh39lf39fU9cO0ysff2V8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAtRtmALUbZge",
		},
		{
			Name: "ab_sr",
			Value: "1.0.0_OTU4MWU4MWI5ZDJhYzFlMjYxMzJkNmFhN2M5YWRmMjFmODk2NzdkOWU0NzExM2Q5OGY4YTc0MzFjNGY5Mjg1YWVkMzlmNTcwMjFjOTVmZmQyMTcyYzNlNzE2YjRmZWRm",
		},
		//备用的新cookie
		//{
		//	Name: "BDUSS",
		//	Value: "",
		//},
	}

	return cookies
}

// 解析百度指标的搜索关键词
func (b baidu) ParseWords(u string) string {
	words := ""
	reg := regexp.MustCompile(`(\?|&)words=.*?(&|$)`)
	words = reg.FindString(u)
	words = strings.Trim(words, `?words=&`)
	return words
}

// 解析获得百度指标的解密密钥ptbk
func (b baidu) ParsePtbk(u string) string {

	apiHost := b.ApiHosts.IndexHost
	// 经过观察以下参数是固定值
	uniqid := u

	// 格式化获取解密密钥ptbk的链接
	//http://index.baidu.com/Interface/ptbk?uniqid=e24e84bc88aa10a7486f842c850f2b1b
	u = fmt.Sprintf("%s/Interface/ptbk?uniqid=%s",
		apiHost, // 网站host
		uniqid,
	)
	c := b.Collector("详情页")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false
	ptbk := ""

	c.OnResponse(func(r *colly.Response) {

		findFlag = true
		ctxStr := string(r.Body)
		ctxByte := []byte(ctxStr)
		ptbk, _ = jsonparser.GetString(ctxByte, "data")
	})

	_ = c.Visit(u)
	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到解密密钥ptbk", u))
	}

	return ptbk
}

// 解密百度指标的加密数据
// 格式类似于encryptData = 'A4P4tAyo9'，ptbk = 'OGi4PtoBAy83sL9%2.63,9847-15+0'
func (b baidu) DecryptData(key string, data string) string {
	decryptData := ""
	tempKeyMap := make(map[string]string)
	halfKeyLength := int64(len(key) / 2)
	startSlice := key[halfKeyLength:]
	endSlice := key[:halfKeyLength]
	for i := 0; i < len(startSlice); i++ {
		tempKeyMap[string(endSlice[i])] = string(startSlice[i])
	}
	for _, ch := range data {
		decryptData += tempKeyMap[string(ch)]
	}

	return decryptData
}

// 根据返回结果集，返回最近一天百度指标数
func (b baidu) ParseLastDayData(data string) int64 {
	dataParse := strings.Split(data, ",")
	lastDayData := dataParse[len(dataParse)-1]
	lastDayDataInt, _ := strconv.ParseInt(lastDayData, 10, 64)

	return lastDayDataInt
}
