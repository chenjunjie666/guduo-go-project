package comment_count

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func bilibiliHandle() {
	urls := storage.Bilibili.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go bilibiliCommentCount(row.Url, row.ShowId)
	}

	wg.Done()
}

func bilibiliCommentCount(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	epids := common.Bilibili.ParseEpids(u)

	// 评论总数
	tcc := int64(0)
	// 成功获取评论的集数
	succCount := 0

	// 获取单集的评论数函数
	epCommentCount := func(epid string) {
		c := common.Bilibili.Collector(ModName)

		// 如果解析 ep id 失败，记录错误，直接返回
		if epid == "" {
			log.Warn(fmt.Sprintf("未获取到%s的ep——id，源网址为：%s", c.Info.Name, u))
			return
		}

		// 构建获取评论数的 api 的 url
		apiUrl := fmt.Sprintf("%s/pgc/season/episode/web/info?ep_id=%s", common.Bilibili.ApiHosts.ApiHost, epid)

		c.OnError(func(r *colly.Response, e error) {
			c.Retry(r, e)
		})

		c.OnResponse(func(r *colly.Response) {
			jsonByte := r.Body

			commentCount, e := jsonparser.GetInt(jsonByte, "data", "stat", "reply")
			if e != nil {
				log.Warn(fmt.Sprintf("解析哔哩哔哩评论失败：%s，原数据为：%s", e, string(jsonByte)))
				return
			}

			tcc += commentCount
			succCount++
		})

		_ = c.Visit(apiUrl)
	}

	for _, epid := range epids {
		epCommentCount(epid)
	}

	if succCount != len(epids) {
		log.Warn(fmt.Sprintf("获取部分总评论数失败，成功获取:%d个集数，一共需要获取%d个，链接：%s", succCount, len(epids), u))
	}

	// 只有所有集数的评论都完成了收集，才存数据库
	storage.Bilibili.StoreCommentCount(tcc, JobAt, showId)
}
