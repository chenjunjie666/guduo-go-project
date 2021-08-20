package base_info

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestIqiyiIntroduce(t *testing.T) {
	core.Init()

	wg.Add(1)
	row := storage.Iqiyi.GetNeedFetchBaseInfoUrl()[0]
	iqiyiIntroduction(row.Url, row.ShowId)
	// iqiyiIntroduction("https://www.iqiyi.com/v_19rrcuh1jw.html", 100)    // 电影
	// iqiyiIntroduction("https://www.iqiyi.com/v_wu9wftmpq4.html", 100)    // 电视剧
	// iqiyiIntroduction("https://www.iqiyi.com/v_1vqcunfvcuw.html", 100)   // 综艺(没有演员、导演信息)
	// iqiyiIntroduction("https://www.iqiyi.com/v_1whaybd4dn4.html", 100)   // 动漫(没有演员、导演信息)
}
