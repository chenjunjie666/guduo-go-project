package base_info

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strings"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func tencentHandle() {
	urls := storage.Tencent.GetNeedFetchBaseInfoUrl()

	wg.Add(len(urls))
	for _, row := range urls {
		go tencentIntroduction(row.Url, row.ShowId)
	}

	wg.Done()
}

// 爬取演员列表
func tencentIntroduction(u string, showId uint64) {
	defer wg.Done()
	c := common.Tencent.Collector(ModName)

	findFlag := false
	c.OnError(func(r *colly.Response, err error) {
		c.Retry(r, err)
	})

	var baseInfoMap map[string]string
	baseInfoMap = make(map[string]string)

	// 匹配 css 规则获取演员列表内容
	c.OnHTML("div.director", func(ele *colly.HTMLElement) {
		findFlag = true
		baseInfo := ele.Text
		var name []string
		ele.ForEach("div > a", func(_ int, eleItem *colly.HTMLElement) {
			name = append(name, eleItem.Text)
		})
		//return
		for _, v := range name {
			nameIndex := strings.Index(baseInfo, v)
			directorIndex := strings.Index(baseInfo, "导演:")
			actorIndex := strings.Index(baseInfo, "演员:")
			guestIndex := strings.Index(baseInfo, "嘉宾:")

			if nameIndex > directorIndex && (nameIndex < actorIndex || nameIndex < guestIndex) {
				if baseInfoMap["Director"] == "" {
					baseInfoMap["Director"] = v
				} else {
					baseInfoMap["Director"] = baseInfoMap["Director"] + " , " + v
				}
			} else if actorIndex >= 0 && nameIndex > actorIndex {
				if baseInfoMap["Actor"] == "" {
					baseInfoMap["Actor"] = v
				} else {
					baseInfoMap["Actor"] = baseInfoMap["Actor"] + " , " + v
				}
			} else if guestIndex >= 0 && nameIndex > guestIndex {
				if baseInfoMap["Guest"] == "" {
					baseInfoMap["Guest"] = v
				} else {
					baseInfoMap["Guest"] = baseInfoMap["Guest"] + " , " + v
				}
			} else if directorIndex == -1 && actorIndex == -1 {
				fmt.Println("未找到艺人信息")
			}
		}

		// 存储获取到的简介
		storage.Tencent.StoreBaseInfoMap(baseInfoMap, showId)
	})

	_ = c.Visit(u)

	// 如果没有找到，记录错误日志
	if !findFlag {
		log.Warn(fmt.Sprintf("获取演员列表失败，链接：%s", u))
	}

}
