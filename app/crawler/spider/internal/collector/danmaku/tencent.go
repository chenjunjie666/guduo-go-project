package danmaku

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"sync"
	"time"

	"github.com/buger/jsonparser"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func tencentHandle() {

	detailUrls := storage.Tencent.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go tencentDanmakuContent(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取腾讯视频弹幕内容
func tencentDanmakuContent(u string, sid uint64) {
	defer wg.Done()
	coverId := common.Tencent.ParseCoverId(u)

	c := common.Tencent.Collector("腾讯" + ModName)
	c.Async = true
	unqionWorkerId := fmt.Sprintf("%s_%d", "tencent_", time.Now().UnixNano())

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	// 以下为常数值
	otype := "json"
	timeStamp := int64(0)
	vids := common.Tencent.ParseVids(u)
	targetIds := common.Tencent.ParseTargetIds(vids)

	danmakuSlice := make([]string, 0, 10000)
	danmakuIdSlice := make([]string, 0, 10000)
	var lastGet []byte
	errCount := 0

	tc := int64(0)
	mapLock := &sync.Mutex{}
	c.OnResponse(func(r *colly.Response) {
		// 上次拿到的和这次拿到的结果一样，则认为已经结束了
		if string(lastGet) == string(r.Body) {
			errCount++
			return
		} else {
			lastGet = r.Body
			errCount = 0
		}
		// 如果没有获取到commentid则表示这个阶段为空，空计数+1
		// 如果连续空计数超过120次（30分钟）就认为结束了
		if errCount > 120 {
			return
		}
		// 完结校验结束
		dcTemp := r.Body
		jsonparser.ArrayEach(dcTemp, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			tc++
			findFlag = true
			mapLock.Lock()
			id, _ := jsonparser.GetString(value, "commentid")
			ct, _ := jsonparser.GetString(value, "content")
			danmakuSlice = append(danmakuSlice, ct)
			danmakuIdSlice = append(danmakuIdSlice, id)
			mapLock.Unlock()
		}, "comments")
	})

	VideoIndex := 120
	if len(targetIds) <= 2 {
		VideoIndex = 400
	}

	trc := 0
	for _, targetArr := range targetIds {
		vId := targetArr[0]
		targetId := targetArr[1]
		second := common.Tencent.ParseLength(coverId, vId)
		VideoIndex = second/30 + 1
		trc += VideoIndex
		//for ; timeStamp <= 3*60*60; timeStamp += 15 {
		idx := 0
		for idx <= VideoIndex {
			url := fmt.Sprintf("%s/danmu?otype=%s&target_id=%s&timestamp=%d",
				common.Tencent.ApiHosts.VideoMFMHost,
				otype,
				targetId,
				timeStamp,
			)
			_ = c.Visit(url, unqionWorkerId)

			timeStamp += 30
			idx++
		}
	}

	c.Wait()
	c.ReleaseVip(unqionWorkerId)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕", u))
		//return
	}
	// 不要用defer，应为等待数据库操作是很费事的
	ch.PopJob()

	log.Info("++++++++++++++++++++++++腾讯视频弹幕，showid", sid, " 抓取完毕，共抓取", len(targetIds), "集，总共", tc, "条弹幕，准备开始存储++++++++++++++++++++++++")
	// 存储获取到的弹幕
	storage.Tencent.StoreDanmakuContent(danmakuSlice, danmakuIdSlice, JobAt, sid)
	storage.Tencent.StoreDanmakuCount(tc, JobAt, sid)
}
