package indicator_actor

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestQihooIndicator(t *testing.T) {
	core.Init()

	wg.Add(3)
	url := storage.Qihoo.GetIndicatorUrl()[0]
	qihooIndicator(url.Url, url.ShowId)
}
