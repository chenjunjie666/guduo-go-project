package hot

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestYoukuHot(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Youku.GetDetailUrl()[0]
	youkuHot(url.Url, url.ShowId)
	// youkuHot("https://v.youku.com/v_show/id_XNTU1MTI2NTky.html", 100)         // 电视剧
	// youkuHot("https://v.youku.com/v_show/id_XNTEzNTgwOTUyOA==.html", 100)     // 电影
	// youkuHot("https://v.youku.com/v_show/id_XNTEyOTg4NDgyOA==.html", 100)     // 综艺
	// youkuHot("https://v.youku.com/v_show/id_XNTQwMTgxMTE2.html", 100)         // 动漫
}
