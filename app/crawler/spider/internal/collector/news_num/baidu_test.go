package news_num

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestBaiduNewsNum(t *testing.T) {
	core.Init()

	wg.Add(1)
	ch.PushJob()
	//url := storage.Baidu.GetNewsUrl()[0]
	//baiduNewsNum(url.Url, url.ShowId)
	baiduNewsNum("https://www.baidu.com/s?rtt=1&bsst=1&cl=undefined&tn=news&rsv_dl=ns_pc&word=司藤", 0)
}
