package proxy

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/pkg/time"
	time2 "time"
)

var lastFetch int64 = 0

func fetchProxyUrls () proxyUrls {
	t := time2.Now().Unix()
	if t - lastFetch < 4 {
		time2.Sleep(time2.Second * 2)
	}
	lastFetch = t

	ps := make(proxyUrls, 0, 50)

	apiUrl := "http://webapi.http.zhimacangku.com/getip?num=20&type=2&pro=0&city=0&yys=0&port=11&pack=98653&ts=1&ys=0&cs=0&lb=1&sb=0&pb=45&mr=1&regions="
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		code, e := jsonparser.GetInt(r.Body, "code")
		if e != nil || code != 0 {
			log.Error(fmt.Sprintf("获取代理失败，返回code为：%d, 错误是：%s", code, e))
			return
		}

		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			ip, _ := jsonparser.GetString(v, "ip")
			port, _ := jsonparser.GetInt(v, "port")
			expiredStr, _ := jsonparser.GetString(v, "expire_time")
			expired := time.YYYYmmddToSecTimestamp(expiredStr)

			url := fmt.Sprintf("http://%s:%d", ip, port)
			row, _ := newProxyUrl(url, int64(expired))
			ps = append(ps, row)
		}, "data")
	})

	c.Visit(apiUrl)
	return ps
}
