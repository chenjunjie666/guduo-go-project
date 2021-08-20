package test

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core/proxy"
	"testing"
)

func TestRetry(t *testing.T) {
	proxy.InitProxyPool()
	c := common.Youku.CollectorWithToken("test")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("aaa")
	})
}