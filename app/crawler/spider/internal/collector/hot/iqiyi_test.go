package hot

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestIqiyiHot(t *testing.T) {
	core.Init()

	wg.Add(1)
	ch.PushJob()
	//url := storage.Iqiyi.GetDetailUrl()[0]
	url := "https://www.iqiyi.com/v_fuudooxv8k.html"
	iqiyiHot(url, 0)
	// iqiyiHot("https://www.iqiyi.com/v_19rrcuh1jw.html", 100)    // 电影
	// iqiyiHot("https://www.iqiyi.com/v_wu9wftmpq4.html", 100)    // 电视剧
	// iqiyiHot("https://www.iqiyi.com/v_1vqcunfvcuw.html", 100)   // 综艺
	// iqiyiHot("https://www.iqiyi.com/v_1whaybd4dn4.html", 100)   // 动漫
}
