package base_info

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	guduoJson "guduo/pkg/json"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func iqiyiHandle() {
	urls := storage.Iqiyi.GetNeedFetchBaseInfoUrl()

	wg.Add(len(urls))
	for _, row := range urls {
		go iqiyiIntroduction(row.Url, row.ShowId)
	}

	wg.Done()
}

// 爬取演员列表
func iqiyiIntroduction(u string, showId uint64) {
	defer wg.Done()
	c := common.Iqiyi.Collector(ModName)

	findFlag := false
	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	var baseInfoMap map[string]string
	baseInfoMap = make(map[string]string)
	var playRoleMap map[string]string
	playRoleMap = make(map[string]string)
	// 匹配 css 规则获取演员列表内容
	c.OnHTML("ul.intro-detail", func(ele *colly.HTMLElement) {
		findFlag = true
		ele.ForEach(".name-wrap", func(_ int, eleItem *colly.HTMLElement) {
			itemProp := eleItem.ChildAttr("a", "itemprop")
			name := eleItem.ChildAttr("a", "title")
			if itemProp == "director" {
				if baseInfoMap["Director"] == "" {
					baseInfoMap["Director"] = name
				} else {
					baseInfoMap["Director"] = baseInfoMap["Director"] + " , " + name
				}
			} else if itemProp == "actor" {
				if baseInfoMap["Actor"] == "" {
					baseInfoMap["Actor"] = name
					playRoleMap[name] = eleItem.ChildText("span")
				} else {
					baseInfoMap["Actor"] = baseInfoMap["Actor"] + " , " + name
					playRoleMap[name] = eleItem.ChildText("span")
				}
			} else {
				fmt.Println("未找到艺人信息")
			}
		})
		playRoleJsonStr, _ := guduoJson.ConvertToJsonStr(playRoleMap)
		baseInfoMap["PlayRole"] = playRoleJsonStr
		// 存储获取到的简介
		storage.Iqiyi.StoreBaseInfoMap(baseInfoMap, showId)
	})

	_ = c.Visit(u)

	// 如果没有找到，记录错误日志
	if !findFlag {
		log.Warn(fmt.Sprintf("获取演员列表失败，链接：%s", u))
	}
}
