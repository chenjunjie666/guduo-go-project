package comment_count

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"strconv"
	"time"
)

func souhuHandle() {
	urls := storage.Souhu.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go souhuCommentCount(row.Url, row.ShowId)
	}
	wg.Done()
}

func souhuCommentCount(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	vids := common.Souhu.ParseVids(u)
	sc := 0
	tcc := int64(0)

	eachCommentCount := func(vid string) {
		c := common.Souhu.Collector(ModName)

		topicID := vid
		topicType := "1"
		source := "2"
		pageSize := "10"
		sort := "0"
		timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		ssl := "0"
		pageNo := "1"
		replySize := "2"
		callback := common.Souhu.GenJQueryCallbackStr()
		_fd := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
		apiUrl := fmt.Sprintf("%s/comment/h5/load?"+
			"topic_id=%s&topic_type=%s&source=%s&page_size=%s&sort=%s&timestamp=%s"+
			"&ssl=%s&page_no=%s&reply_size=%s&callback=%s&_=%s",
			common.Souhu.ApiHosts.ApiHost,
			topicID,
			topicType,
			source,
			pageSize,
			sort,
			timestamp,
			ssl,
			pageNo,
			replySize,
			callback,
			_fd,
		)

		c.OnResponse(func(r *colly.Response) {
			j := string(r.Body)

			reg := regexp.MustCompile(`comment_count":\d+`)
			ccTmp := reg.FindString(j)

			reg = regexp.MustCompile(`\d+`)
			ccStr := reg.FindString(ccTmp)
			cc, e := strconv.ParseInt(ccStr, 10, 64)

			if e != nil {
				log.Warn(fmt.Sprintf("转换评论数失败：%s，源数据：%s", e, j))
				return
			}
			sc++
			tcc += cc
		})

		_ = c.Visit(apiUrl)
	}

	for _, vid := range vids {
		eachCommentCount(vid)
	}

	if sc != len(vids) {
		log.Warn(fmt.Sprintf("获取部分总评论数失败，成功获取：%d，总共:%d，链接：%s", sc, len(vids), u))
	}
	log.Info("++++++++++++++++++获取 搜狐视频 评论结束，showid：", showId, "一共", len(vids), "集, 成功爬取", sc, "集, 一共爬取到：", tcc, "条评论")

	storage.Souhu.StoreCommentCount(tcc, JobAt, showId)
}
