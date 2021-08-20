package indicator_actor

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestWeiboIndicator(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Weibo.GetIndicatorUrl()[0]
	weiboIndicator(url.Name, url.Id)
}
