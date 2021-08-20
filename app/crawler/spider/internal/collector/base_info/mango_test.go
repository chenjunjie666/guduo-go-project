package base_info

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestMangoIntroduce(t *testing.T) {
	core.Init()

	wg.Add(1)
	row := storage.Mango.GetNeedFetchBaseInfoUrl()[0]
	mangoIntroduction(row.Url, row.ShowId)
	// mangoIntroduction("https://www.mgtv.com/b/336049/11619503.html", 100)     // 电影
	// mangoIntroduction("https://www.mgtv.com/b/342104/11630288.html", 100)     // 电视剧
	// mangoIntroduction("https://www.mgtv.com/b/294264/3286815.html", 100)      // 动漫
	// mangoIntroduction("https://www.mgtv.com/b/368214/11859385.html", 100)     // 综艺(没有演员、导演信息)
}
