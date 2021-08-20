package hot

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func iqiyiHandle() {
	urls := storage.Iqiyi.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go iqiyiHot(row.Url, row.ShowId)
	}
	wg.Done()
}

// 爱奇艺热度趋势爬取主逻辑
func iqiyiHot(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	vid := common.Iqiyi.ParseVid(u) // 根据详情页链接，解析得到视频的唯一ID

	if vid == "" {
		log.Warn("未获取到爱奇艺热度趋势api所需要的id参数")
		return
	}

	// 根据vid构建获取热度的api链接
	apiUrl := fmt.Sprintf("https://pcw-api.iqiyi.com/video/video/hotplaytimes/%s", vid)

	c := common.Iqiyi.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false
	// 分析返回值
	c.OnResponse(func(r *colly.Response) {
		jsonStr := r.Body // 返回内容为json
		// hot = json["data"]["hot"]
		_, _ = jsonparser.ArrayEach(r.Body, func(vv []byte, dataType jsonparser.ValueType, offset int, err error) {
			id, _ := jsonparser.GetString(vv, "id")
			if id == vid {
				hot, e := jsonparser.GetInt(vv, "hot")
				if e != nil {
					log.Warn(fmt.Sprintf("解析爱奇艺热度趋势api返回的json失败，原因：%s，源数据：%s", e.Error(), jsonStr))
					return
				}
				findFlag = true
				// 存储爱奇艺热度趋势
				storage.Iqiyi.StoreHot(hot, JobAt, showId)
				return
			}
		}, "data")
	})

	_ = c.Visit(apiUrl)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到热度趋势", u))
	}
}
