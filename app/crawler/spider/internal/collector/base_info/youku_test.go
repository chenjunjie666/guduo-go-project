package base_info

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestYoukuIntroduce(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Youku.GetNeedFetchBaseInfoUrl()[0]
	youkuIntroduction(url.Url, url.ShowId)
	// youkuIntroduction("https://v.youku.com/v_show/id_XNTU1MTI2NTky.html", 100)         // 电视剧
	// youkuIntroduction("https://v.youku.com/v_show/id_XNTEzNTgwOTUyOA==.html", 100)     // 电影
	// youkuIntroduction("https://v.youku.com/v_show/id_XNTEyOTg4NDgyOA==.html", 100)     // 综艺(没有演员、导演信息)
	// youkuIntroduction("https://v.youku.com/v_show/id_XNTQwMTgxMTE2.html", 100)         // 动漫(没有演员、导演信息)
}
