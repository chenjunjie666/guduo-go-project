package common

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/gocolly/colly/v2/extensions"
)

var Douban = &douban{
	PlatformId: storage.Douban.PlatformId,
	Host:       storage.Douban.Host,
}

type douban struct {
	PlatformId uint64
	Host       string
}

// 豆瓣电影采集器初始化
func (d douban) Collector(mod string) *core.CollectorObj {
	cInfo := &core.CollectorInfo{
		"豆瓣电影",
		d.Host,
		d.PlatformId,
		mod,
	}

	// 爬虫配置
	c := core.NewCollector(cInfo) // 初始化爬虫

	extensions.RandomUserAgent(c.Collector) // 随机设置 user-agent
	c.UseProxy()                            // 使用代理，如果代理没有则不使用
	c.DetectCharset = true                  // 非utf-8字符集支持
	return c
}
