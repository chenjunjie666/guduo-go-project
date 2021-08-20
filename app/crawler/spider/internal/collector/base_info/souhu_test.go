package base_info

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestSouhuIntroduce(t *testing.T) {
	core.Init()

	wg.Add(1)
	row := storage.Souhu.GetNeedFetchBaseInfoUrl()[0]
	souhuIntroduction(row.Url, row.ShowId)
	// souhuIntroduction("https://tv.sohu.com/v/MjAyMTA0MjkvbjYwMTAwMzU1MC5zaHRtbA==.html", 100) // 电视剧
	// souhuIntroduction("https://film.sohu.com/album/9702422.html", 100)                        // 电影
	// souhuIntroduction("https://tv.sohu.com/v/MjAyMTA1MDQvbjYwMTAwNDMzNC5zaHRtbA==.html", 100) // 综艺(没有演员、导演信息)
	// souhuIntroduction("https://tv.sohu.com/v/MjAyMTA0MTMvbjYwMDk5NjM4OS5zaHRtbA==.html", 100) // 动漫(没有演员、导演信息)
}
