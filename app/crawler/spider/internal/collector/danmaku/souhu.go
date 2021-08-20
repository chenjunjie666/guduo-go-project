package danmaku

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"math"
	"strings"
	"sync"
	"time"
)

func souhuHandle() {

	detailUrls := storage.Souhu.GetDetailUrl()

	for _, row := range detailUrls {
		if !strings.Contains(row.Url, "com/v/") {
			continue
		}
		wg.Add(1)
		ch.PushJob()
		go souhuDanmakuContent(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取搜狐视频弹幕内容
func souhuDanmakuContent(u string, sid uint64) {
	defer wg.Done()
	c := common.Souhu.Collector("搜狐" + ModName)
	c.Async = true

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	tc := int64(0)
	findFlag := false

	vids := common.Souhu.ParseVids(u)

	if len(vids) == 0 {
		log.Warning("show_id:", sid, "一个vid都没爬到，连接为：", u)
	}

	aid := common.Souhu.ParsePlayListId(u)
	//经观察以下为常数值
	act := "dmlist_v2"
	requestFrom := "h5_js"
	timeBegin := 0
	timeEnd := 100000

	danmakuSlice := make([]string, 0, 5000)
	mapLock := &sync.Mutex{}
	c.OnResponse(func(r *colly.Response) {
		findFlag = true
		dcTemp := r.Body
		jsonparser.ArrayEach(dcTemp, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			tc++
			createdTimeStamp, _ := jsonparser.GetFloat(value, "dd")
			if filterLastDanmaku(int64(JobAt), int64(math.Floor(createdTimeStamp))) {
				mapLock.Lock()
				ct , _ := jsonparser.GetString(value, "c")
				danmakuSlice = append(danmakuSlice, ct)
				mapLock.Unlock()
			}
		}, "info", "comments")
	})
	unqionWorkerId := fmt.Sprintf("%s_%d", "mgtv_", time.Now().UnixNano())
	for _, vid := range vids {
		url := fmt.Sprintf("%s/dmh5/dmListAll?act=%s&request_from=%s&vid=%s&aid=%s&time_begin=%d&time_end=%d",
			common.Souhu.ApiHosts.DanmuApiHost,
			act,
			requestFrom,
			vid,
			aid,
			timeBegin,
			timeEnd,
		)
		_ = c.Visit(url, unqionWorkerId)
	}

	c.Wait()

	c.ReleaseVip(unqionWorkerId)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕数", u))
		//return
	}
	log.Info("++++++++++++++++++++++++搜狐视频弹幕，showid", sid, " 抓取完毕，共抓取", len(vids), "集，总共", tc, "条弹幕，准备开始存储++++++++++++++++++++++++")
	// 不要用defer，应为等待数据库操作是很费事的
	ch.PopJob()
	storage.Souhu.StoreDanmakuContent(danmakuSlice, JobAt, sid)
	storage.Souhu.StoreDanmakuCount(tc, JobAt, sid)
}
