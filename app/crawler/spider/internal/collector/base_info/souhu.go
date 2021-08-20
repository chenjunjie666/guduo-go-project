package base_info

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strings"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func souhuHandle() {
	urls := storage.Souhu.GetNeedFetchBaseInfoUrl()

	wg.Add(len(urls))
	for _, row := range urls {
		go souhuIntroduction(row.Url, row.ShowId)
	}

	wg.Done()
}

// 爬取演员列表
func souhuIntroduction(u string, showId uint64) {
	defer wg.Done()
	// 从详情页获取视频的vid
	// 搜狐视频的链接格式：
	apiHost := common.Souhu.ApiHosts.PlHdHost
	playlistId := common.Souhu.ParsePlayListId(u)
	o_playlistId := common.Souhu.ParseOPlayListId(u)
	// 经过观察此参数是固定值
	pagesize := "999"
	cbNo := common.Souhu.GenJQueryCallbackStr()

	// 格式化获取演员信息的链接
	u = fmt.Sprintf("%s/videolist?playlistid=%s&o_playlistId=%s&callback=%s&pagesize=%s",
		apiHost,      // 网站host
		playlistId,   // playlistid 参数
		o_playlistId, // o_playlistId 参数
		cbNo,         // callback 参数
		pagesize,     // pagesize 参数
	)

	c := common.Souhu.Collector(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false
	var baseInfoMap map[string]string
	baseInfoMap = make(map[string]string)

	c.OnResponse(func(r *colly.Response) {
		// 返回值格式为 jsonp_xxxx_xxxx(json_str)
		// 所以需要将返回值转为 string
		// 然后去掉 "jsonp_xxxx_xxxx(" 以及 ")" 去掉，得到正确的 json 字符串
		ctxStr := string(r.Body)

		ctxStr = strings.TrimLeft(ctxStr, cbNo+"(")
		ctxStr = strings.TrimRight(ctxStr, ");")
		ctxByte := []byte(ctxStr)

		_, _ = jsonparser.ArrayEach(ctxByte, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			name := string(value)
			if baseInfoMap["Actor"] == "" {
				baseInfoMap["Actor"] = name
			} else {
				baseInfoMap["Actor"] = baseInfoMap["Actor"] + " , " + name
			}
		}, "actors")

		findFlag = true
		storage.Souhu.StoreBaseInfoMap(baseInfoMap, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到播放量", u))
	}
}
