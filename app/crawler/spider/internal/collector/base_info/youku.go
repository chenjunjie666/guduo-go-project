package base_info

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	guduoJson "guduo/pkg/json"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func youkuHandle() {
	urls := storage.Youku.GetNeedFetchBaseInfoUrl()

	wg.Add(len(urls))
	for _, row := range urls {
		go youkuIntroduction(row.Url, row.ShowId)
	}

	wg.Done()
}

// 爬取演员列表
func youkuIntroduction(u string, showId uint64) {
	defer wg.Done()
	c := common.Youku.Collector(ModName)

	findFlag := false
	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	var baseInfoMap map[string]string
	baseInfoMap = make(map[string]string)
	var playRoleMap map[string]string
	playRoleMap = make(map[string]string)
	// 匹配 css 规则获取演员列表内容
	c.OnResponse(func(r *colly.Response) {
		findFlag = true
		ctxByte := common.Youku.ParseHtmlInitDataJson(u)
		_, _ = jsonparser.ArrayEach(ctxByte, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			_, e := jsonparser.GetInt(value, "id")
			name := ""
			subtitle := ""
			if e != nil {
				name, _ = jsonparser.GetString(value, "data", "title")
				subtitle, _ = jsonparser.GetString(value, "data", "subtitle")
			}
			if subtitle == "导演" {
				if baseInfoMap["Director"] == "" {
					baseInfoMap["Director"] = name
				} else {
					baseInfoMap["Director"] = baseInfoMap["Director"] + " , " + name
				}
			} else {
				if baseInfoMap["Actor"] == "" {
					baseInfoMap["Actor"] = name
					playRoleMap[name] = subtitle
				} else {
					baseInfoMap["Actor"] = baseInfoMap["Actor"] + " , " + name
					playRoleMap[name] = subtitle
				}
			}
		}, "data", "data", "nodes", "[0]", "nodes", "[0]", "nodes")
		playRoleJsonStr, _ := guduoJson.ConvertToJsonStr(playRoleMap)
		baseInfoMap["PlayRole"] = playRoleJsonStr

		// 存储获取到的简介
		storage.Youku.StoreBaseInfoMap(baseInfoMap, showId)
	})

	_ = c.Visit(u)

	// 如果没有找到，记录错误日志
	if !findFlag {
		log.Warn(fmt.Sprintf("获取演员列表失败，链接：%s", u))
	}

}
