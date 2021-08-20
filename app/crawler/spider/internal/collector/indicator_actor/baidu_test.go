package indicator_actor

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestBaiduIndicator(t *testing.T) {
	core.Init()

	wg.Add(3)
	url := storage.Baidu.GetIndicatorUrl()[0]
	//baiduGenderIndicator(url.IndicatorUrl, url.ShowId)
	//baiduAgeIndicator(url.GenAgeIndicatorUrl, url.ShowId)
	baiduIndicator(url.IndicatorUrl, url.ShowId)
}
