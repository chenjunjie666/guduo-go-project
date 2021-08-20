package base_info

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestTencentIntroduce(t *testing.T) {
	core.Init()
	wg.Add(1)
	row := storage.Tencent.GetNeedFetchBaseInfoUrl()[0]
	tencentIntroduction(row.Url, row.ShowId)
	// tencentIntroduction("https://v.qq.com/x/cover/zf2z0xpqcculhcz.html", 100) // 电视剧
	// tencentIntroduction("https://v.qq.com/x/cover/mzc00200yuiy84o.html", 100) // 电影
	// tencentIntroduction("https://v.qq.com/x/cover/hwm1ryf35f35wed.html", 100) // 综艺
	// tencentIntroduction("https://v.qq.com/x/cover/m441e3rjq9kwpsc.html", 100) // 动漫(没有演员、导演信息)
}
