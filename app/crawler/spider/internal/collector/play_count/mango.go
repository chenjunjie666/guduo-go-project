package play_count

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func mangoHandle() {
	detailUrls := storage.Mango.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go mangoPlayCount(row.Url, row.ShowId)
	}

	wg.Done()
}

// 经过观察某些视频没有显示播放量，故抓取不到
func mangoPlayCount(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()

	ou := u

	// 从详情页获取视频的vid
	// 芒果TV的链接格式：
	// https://www.mgtv.com/.../{vid}.html?...
	vids := common.Mango.ParseVids(u)
	if len(vids) == 0 {
		vids = []string{common.Mango.ParseVid(u)}
	}

	c := common.Mango.Collector(ModName)

	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	findFlag := false
	cbNo := common.Mango.GenJsonpCallbackStr()

	tpc := int64(0)
	c.OnResponse(func(r *colly.Response) {
		// 返回值格式为 jsonp_xxxx_xxxx(json_str)
		// 所以需要将返回值转为 string
		// 然后去掉 "jsonp_xxxx_xxxx(" 以及 ")" 去掉，得到正确的 json 字符串
		ctxStr := string(r.Body)
		ctxStr = strings.TrimLeft(ctxStr, cbNo+"(")
		ctxStr = strings.TrimRight(ctxStr, ")")
		ctxByte := []byte(ctxStr)

		// pc = json["data"]["all"]
		pc, err := jsonparser.GetInt(ctxByte, "data", "all")
		if err != nil {
			log.Warn(fmt.Sprintf("解析json失败,原json为：%s，失败原因为：%s", ctxStr, err))
			return
		}

		findFlag = true
		tpc += pc
	})

	apiHost := common.Mango.ApiHosts.VcHost
	// 经过观察下面四个是固定值
	spt := "10000000"
	abd := "undefined"
	vvType := "1"
	_type := "4"
	for _, vid := range vids {
		// 格式化获取播放量的链接
		apiUrl := fmt.Sprintf("%s/v2/dynamicinfo?_support=%s&vid=%s&abroad=%s&vvType=%s&type=%s&callback=%s",
			apiHost, // 网站host
			spt,     // _support 参数
			vid,     // vid 参数
			abd,     // abroad 参数
			vvType,  // vvType 参数
			_type,   // type 参数
			cbNo,    // callback 参数
		)
		_ = c.Visit(apiUrl)
	}

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到播放量", ou))
	}
	storage.Mango.StorePlayCount(tpc, JobAt, sid)
}
