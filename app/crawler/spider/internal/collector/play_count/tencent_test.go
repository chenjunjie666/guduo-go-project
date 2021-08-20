package play_count

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestTencentPlayCount(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Tencent.GetDetailUrl()[0]
	tencentPlayCount(url.Url, url.ShowId)
	// tencentPlayCount("https://v.qq.com/x/cover/zf2z0xpqcculhcz.html", 100) // 电视剧
	// tencentPlayCount("https://v.qq.com/x/cover/mzc00200yuiy84o.html", 100) // 电影
	// tencentPlayCount("https://v.qq.com/x/cover/hwm1ryf35f35wed.html", 100) // 综艺
	// tencentPlayCount("https://v.qq.com/x/cover/m441e3rjq9kwpsc.html", 100) // 动漫
}
