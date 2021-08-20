package danmaku

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/lib/youku"
	"guduo/app/crawler/spider/internal/storage"
	"time"
)

func youkuHandle() {
	detailUrls := storage.Youku.GetDetailUrl()

	for _, row := range detailUrls {
		// 不是优酷播放页，直接跳过，不爬了
		if !common.Youku.CheckIsYoukuPlayPage(row.Url) {
			continue
		}

		wg.Add(1)
		ch.PushJob()
		go youkuDanmakuContent(row.Url, row.ShowId)
	}

	wg.Done()
}

func youkuDanmakuContent(u string, sid uint64) {
	defer wg.Done()
	unqionWorkerId := fmt.Sprintf("%s_%d", "youku_", time.Now().UnixNano())
	vids := common.Youku.ParseVids(u)
	mins := common.Youku.ParseLength(u)
	if mins > 240 {
		mins = 240
	}
	tc := int64(0)
	dmks := make([]string, 0, 5000)

	c := common.Youku.CollectorWithToken("优酷" + ModName)
	c.Async = true

	c.OnError(func(r *colly.Response, e error) {
		log.Warn("优酷弹幕访问出错：show id：", sid, "错误：", e.Error())
		//c.Retry(r, e)
	})
	startTime := int64((JobAt - 3600*4) * 1000)
	endTime := int64(JobAt * 1000)

	c.OnResponse(func(r *colly.Response) {
		s, e := jsonparser.GetString(r.Body, "data", "result")
		if e != nil {
			api := r.Request.URL.String()
			log.Warn(fmt.Sprintf("解析优酷弹幕json文件外层出错：url:%s", api))
		}

		_, e = jsonparser.ArrayEach([]byte(s), func(v []byte, dataType jsonparser.ValueType, offset int, err error) {
			cTime, _ := jsonparser.GetInt(v, "createtime")
			ctx, _ := jsonparser.GetString(v, "content")
			tc++
			//fmt.Println(startTime, cTime, endTime)
			if cTime >= startTime && cTime < endTime {
				dmks = append(dmks, ctx)
			}
		}, "data", "result")

		if e != nil {
			log.Warn(fmt.Sprintf("解析优酷弹幕json文件内层出错：%s, url:%s 源数据：%s,", e, u, string(s)))
		}
	})

	for _, vid := range vids {
		for i := 0; i <= mins; i++ {
			apiUrl, reqData := youku.GetDanmakuUrl(vid, i)
			c.Post(apiUrl, reqData, unqionWorkerId)
		}
	}

	c.Wait()

	core.ReleaseVip(unqionWorkerId, common.Youku.PlatformId)
	log.Info("++++++++++++++++++++++++优酷视频弹幕，showid", sid, " 抓取完毕，共抓取", len(vids), "集，总共", tc, "条弹幕，准备开始存储++++++++++++++++++++++++")
	if len(dmks) == 0 {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕数", u))
		//return
	}
	// 不要用defer，应为等待数据库操作是很费事的
	ch.PopJob()
	storage.Youku.StoreDanmankuContent(dmks, JobAt, sid)
	storage.Youku.StoreDanmankuCount(tc, JobAt, sid)
}
