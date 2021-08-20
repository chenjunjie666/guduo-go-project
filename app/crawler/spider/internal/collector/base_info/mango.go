package base_info

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	guduoJson "guduo/pkg/json"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func mangoHandle() {
	urls := storage.Mango.GetNeedFetchBaseInfoUrl()

	wg.Add(len(urls))
	for _, row := range urls {
		go mangoIntroduction(row.Url, row.ShowId)
	}
	wg.Done()
}

// 爬取演员列表
func mangoIntroduction(u string, showId uint64) {
	defer wg.Done()
	// 从详情页获取视频的vid
	// 芒果TV的链接格式：
	// https://www.mgtv.com/{cid}/{vid}.html?...
	apiHost := common.Mango.ApiHosts.PcWebHost
	vid := common.Mango.ParseVid(u)
	cid := common.Mango.ParseCid(u)
	// 经过观察此参数是固定值
	spt := "10000000"
	cbNo := common.Mango.GenJsonpCallbackStr()

	// 格式化获取播放量的链接
	u = fmt.Sprintf("%s/video/info?&vid=%s&cid=%s&_support=%s&callback=%s",
		apiHost, // 网站host
		vid,     // vid 参数
		spt,     // _support 参数
		cid,     // cid 参数
		cbNo,    // callback 参数
	)

	c := common.Mango.Collector(ModName)

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
		ctxStr = strings.TrimRight(ctxStr, ")")
		ctxByte := []byte(ctxStr)

		directorTemp := guduoJson.GetFieldFromJson("data/info/detail/director", ctxByte)
		director := strings.Replace(directorTemp, " / ", ",", -1)
		actorTemp := guduoJson.GetFieldFromJson("data/info/detail/leader", ctxByte)
		actor := strings.Replace(actorTemp, " / ", ",", -1)
		baseInfoMap["Director"] = director
		baseInfoMap["Actor"] = actor

		findFlag = true
		storage.Mango.StoreBaseInfoMap(baseInfoMap, showId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到播放量", u))
	}
}
