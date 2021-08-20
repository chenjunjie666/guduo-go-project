package comment_count

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func iqiyiHandle() {
	urls := storage.Iqiyi.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go iqiyiCommentCount(row.Url, row.ShowId)
	}
	wg.Done()
}

func iqiyiCommentCount(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	cids := common.Iqiyi.ParseVids(u)
	if len(cids) == 0 {
		// 电影没有这些花里胡哨的接口
		cidOne := common.Iqiyi.ParseVid(u)
		cids = []string{cidOne}
	}


	accessCount := 0
	tcc := int64(0) // 总评论数

	eachCommentCount := func(cid string) {
		agentType := "118"
		agentVer := "9.11.5"
		authCookie := "null"
		businessType := "17"
		contentId := cid
		hotSize := "10"
		lastId := ""
		page := ""
		pageSize := "10"
		types := "hot,time"
		callback := common.Iqiyi.GenJsonpCallbackStr()

		apiUrl := fmt.Sprintf("%s/v3/comment/get_comments.action?"+
			"agent_type=%s&agent_version=%s&authcookie=%s&business_type=%s&content_id=%s&hot_size=%s&last_id=%s"+
			"&page=%s&page_size=%s&types=%s&callback=%s",
			common.Iqiyi.ApiHosts.SnsHost,
			agentType,
			agentVer,
			authCookie,
			businessType,
			contentId,
			hotSize,
			lastId,
			page,
			pageSize,
			types,
			callback,
		)
		c := common.Iqiyi.Collector(ModName)

		c.OnError(func(r *colly.Response, e error) {
			c.Retry(r, e)
		})

		c.OnResponse(func(r *colly.Response) {
			b := r.Body

			reg := regexp.MustCompile(`"commentReplyCount":\d+`)
			ccTmp := reg.Find(b)

			reg = regexp.MustCompile(`\d+`)
			ccStr := string(reg.Find(ccTmp))

			cc, e := strconv.ParseInt(ccStr, 10, 64)

			if e != nil {
				log.Warn(fmt.Sprintf("爱奇艺评论数转换失败：%s，源数据：%s", e, string(b)))
				return
			}
			accessCount++
			fmt.Println(accessCount, cc)
			tcc += cc
		})

		_ = c.Visit(apiUrl)
	}

	for _, cid := range cids {
		eachCommentCount(cid)
	}

	if accessCount != len(cids) {
		log.Warn(fmt.Sprintf("获取爱奇艺部分评论失败，成功获取：%d，一共：%d", accessCount, len(cids)))
	}
	log.Info("++++++++++++++++++获取爱奇艺评论结束，showid：", showId, "一共", len(cids), "集, 成功爬取", accessCount, "集, 一共爬取到：", tcc, "条评论")
	return
	storage.Iqiyi.StoreCommentCount(tcc, JobAt, showId)
}
