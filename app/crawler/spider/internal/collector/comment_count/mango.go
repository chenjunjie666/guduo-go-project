package comment_count

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func mangoHandle() {
	urls := storage.Mango.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go mangoCommentCount(row.Url, row.ShowId)
	}
	wg.Done()
}

func mangoCommentCount(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	vids := common.Mango.ParseVids(u)

	tcc := int64(0)
	accessCount := 0

	eachCommentCount := func(vid string) {
		apiHost := common.Mango.ApiHosts.CommentHost
		subType := "hunantv2014"
		subId := vid
		uuid := ""
		ticket := ""
		platform := "pc"
		callback := common.Mango.GenJsonpCallbackStr()
		_support := "10000000"
		_fd := time.Now().UnixNano() / 1e6
		apiUrl := fmt.Sprintf("%s/v4/comment/getCount?subjectType=%s&subjectId=%s&uuid=%s&ticket=%s"+
			"&platform=%s&callback=%s&_support=%s&_=%d",
			apiHost,
			subType,
			subId,
			uuid,
			ticket,
			platform,
			callback,
			_support,
			_fd,
		)

		c := common.Mango.Collector(ModName)

		c.OnError(func(r *colly.Response, e error) {
			c.Retry(r, e)
		})

		c.OnResponse(func(r *colly.Response) {
			jTmp := string(r.Body)

			j := strings.Trim(jTmp, callback+")")

			cc, e := jsonparser.GetInt([]byte(j), "data", "commentCount")

			if e != nil {
				log.Warn(fmt.Sprintf("获取芒果TV评论数失败：%s，源数据：%s", e, jTmp))
				return
			}

			tcc += cc
			accessCount++
		})

		_ = c.Visit(apiUrl)
	}

	for _, vid := range vids {
		eachCommentCount(vid)
	}

	if accessCount != len(vids) {
		log.Warn(fmt.Sprintf("获取芒果部分评论失败，成功获取：%d，一共：%d", accessCount, len(vids)))
	}

	log.Info("++++++++++++++++++获取芒果TV评论结束，showid：", showId, "一共", len(vids), "集, 成功爬取", accessCount, "集, 一共爬取到：", tcc, "条评论")
	storage.Mango.StoreCommentCount(tcc, JobAt, showId)
}
