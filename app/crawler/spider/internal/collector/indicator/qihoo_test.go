package indicator

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestQihooIndicator(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := fmt.Sprintf("https://trends.so.com/result?query=%s&period=30", "司藤")
	qihooIndicator(url, 0)
}
