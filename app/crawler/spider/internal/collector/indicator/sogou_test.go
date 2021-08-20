package indicator

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestSogouIndicator(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Sogou.GetIndicatorUrl()[0]
	sogouIndicator(url.Url, url.ShowId)
}
