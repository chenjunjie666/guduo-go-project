package common

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/gocolly/colly/v2/extensions"
)

var Wechat = &wechat{
	PlatformId: storage.Wechat.PlatformId,
	Host:       storage.Wechat.Host,
}

type wechat struct {
	PlatformId uint64
	Host       string
}

// 微信文章数采集器初始化
func (w wechat) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"微信",
		w.Host,
		w.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}
