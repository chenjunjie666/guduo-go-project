package danmaku

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestIqiyiDanmakuContent(t *testing.T) {
	core.Init()
	//url := storage.Iqiyi.GetDetailUrl()[1]

	wg.Add(1)
	url := "https://www.iqiyi.com/v_199f24ura6k.html#curid=4540760052391200_a28f39e231e6e48f3e8583df18659802"
	iqiyiDanmakuContent(url, 15112)
}
