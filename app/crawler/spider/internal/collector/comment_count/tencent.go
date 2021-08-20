package comment_count

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strconv"
	"strings"
	"time"
)

func tencentHandle() {
	urls := storage.Tencent.GetDetailUrl()

	for _, row := range urls {
		wg.Add(1)
		ch.PushJob()
		go tencentCommentCount(row.Url, row.ShowId)
	}
	wg.Done()
}

func tencentCommentCount(u string, showId uint64) {
	defer ch.PopJob()
	defer wg.Done()
	vids := common.Tencent.ParseVids(u)
	cids := common.Tencent.ParseCommentIds(vids)
	if len(cids) == 0 {
		log.Warn(fmt.Sprintf("获取评论失败，评论ID未成功获取：%s", u))
		return
	}

	tcc := int64(0)
	sc := 0

	eachCommentCount := func(cid string) {
		c := common.Tencent.Collector(ModName)

		commentId := cid
		cbStr := common.Tencent.GenVarticleCallbackStr(commentId)
		orinum := "10"
		oriorder := "o"
		pageflag := "1"
		cursor := "0"
		scorecursor := "0"
		orirepnum := "2"
		reporder := "0"
		reppageflag := "1"
		source := "132"
		_fd := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)

		apiUrl := fmt.Sprintf("%s/varticle/%s/comment/v2?callback=%s&orinum=%s&oriorder=%s"+
			"&pageflag=%s&cursor=%s&scorecursor=%s&orirepnum=%s&reporder=%s&reppageflag=%s&source=%s&_=%s",
			common.Tencent.ApiHosts.CommentHost,
			commentId,
			cbStr,
			orinum,
			oriorder,
			pageflag,
			cursor,
			scorecursor,
			orirepnum,
			reporder,
			reppageflag,
			source,
			_fd,
		)

		c.OnError(func(r *colly.Response, e error) {
			c.Retry(r, e)
		})

		c.OnResponse(func(r *colly.Response) {
			jOrigin := string(r.Body)
			j := strings.Trim(jOrigin, cbStr+"()")
			ccStr, e := jsonparser.GetString([]byte(j), "data", "targetInfo", "commentnum")

			if e != nil {
				log.Warn(fmt.Sprintf("解析评论json失败：%s，源数据：%s", e, jOrigin))
				return
			}

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

	for _, cid := range cids {
		eachCommentCount(cid)
	}

	if sc != len(cids) {
		log.Warn(fmt.Sprintf("获取部分腾讯评论数失败，成功获取:%d，总共：%d，连接：%s", sc, len(cids), u))
	}
	log.Info("++++++++++++++++++获取腾讯视频评论结束，showid：", showId, "一共", len(vids), "集, 成功爬取", sc, "集, 一共爬取到：", tcc, "条评论")

	storage.Tencent.StoreCommentCount(tcc, JobAt, showId)
}
