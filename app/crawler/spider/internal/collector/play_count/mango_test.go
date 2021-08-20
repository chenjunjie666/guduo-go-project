package play_count

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestMangoPlayCount(t *testing.T) {
	core.Init()

	wg.Add(1)
	//url := storage.Mango.GetDetailUrl()[0]
	//mangoPlayCount(url.Url, url.ShowId)

	mangoPlayCount("https://www.mgtv.com/b/369217/11964491.html", 0)
	// mangoPlayCount("https://www.mgtv.com/b/336049/11619503.html", 100)     // 电影
	// mangoPlayCount("https://www.mgtv.com/b/342104/11630288.html", 100)     // 电视剧
	// mangoPlayCount("https://www.mgtv.com/b/294264/3286815.html", 100)      // 动漫
	// mangoPlayCount("https://www.mgtv.com/b/368214/11859385.html", 100)     // 综艺
}
