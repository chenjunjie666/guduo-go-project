package attention_actor

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"guduo/app/crawler/spider/internal/util"
	"regexp"
	"strings"
)

func baiduHandle() {

	detailUrls := storage.Baidu.GetActorDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go baiduAttention(row.Url, row.ActorId)
	}

	wg.Done()
}

// 从页面解析并获取获取贴吧关注度
func baiduAttention(u string, actorId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Baidu.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false


	// 获取贴吧关注数和帖子数
	c.OnResponse(func(r *colly.Response) {
		// 匹配贴吧关注数
		reg := regexp.MustCompile(`card_menNum">.*?</span>`)
		aTmp := reg.FindString(string(r.Body))
		aTmp = strings.Trim(aTmp, `card_menNum"></span>`)
		a, e := util.EscapeDotInt(aTmp)
		if e != nil {
			log.Warn(fmt.Sprintf("解析百度贴吧关注数失败：%s，源数据:%s", e, aTmp))
			return
		}

		// 匹配贴吧发帖数
		reg = regexp.MustCompile(`card_infoNum">.*?</span>`)
		pTmp := reg.FindString(string(r.Body))
		pTmp = strings.Trim(pTmp, `card_infoNum"></span>`)
		p, e := util.EscapeDotInt(pTmp)
		if e != nil {
			log.Warn(fmt.Sprintf("解析百度帖子数失败：%s，源数据:%s", e, pTmp))
			return
		}

		findFlag = true
		storage.Baidu.StoreAttentionActor(a, p, JobAt, actorId)
	})

	_ = c.Visit(u)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到贴吧关注度", u))
	}
}
