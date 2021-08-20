package comment_count

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/lib/youku"
	"guduo/app/crawler/spider/internal/storage"
)

func youkuHandle() {
	urls := storage.Youku.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go youkuCommentCount(row.Url, row.ShowId)
	}
	wg.Done()
}

func youkuCommentCount(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()

	vids := common.Youku.ParseVids(u)

	tcc := int64(0)
	sc := 0
	for _, vid := range vids {
		apiUrl := youku.GetCommentApiUrl(vid)
		cc, isFind := youkuFetchCommentCount(apiUrl)
		if !isFind {
			log.Warn(fmt.Sprintf("部分优酷评论数未获取到：%s", u))
			continue
		}
		sc++
		tcc += cc
	}
	//fmt.Println(tcc)
	//return
	log.Info("++++++++++++++++++获取芒果TV评论结束，showid：", sid, "一共", len(vids), "集, 成功爬取", sc, "集, 一共爬取到：", tcc, "条评论")
	storage.Youku.StoreCommentCount(tcc, JobAt, sid)
}

func youkuFetchCommentCount(u string) (int64, bool) {
	c := common.Youku.CollectorWithToken(ModName)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("Host", "acs.youku.com")
	})

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	cc := int64(0)
	isFind:= false
	c.OnResponse(func(r *colly.Response) {
		_, _ = jsonparser.ArrayEach(r.Body, func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			type_, _ := jsonparser.GetString(v, "dataSource")
			if type_ == "ALL_COMMENT_DATASOURCE" {
				isFind = true

				cc, _ = jsonparser.GetInt(v, "totalCount")
			}
		}, "data", "data", "modules")
	})

	_ = c.Visit(u)
	return cc, isFind
}
