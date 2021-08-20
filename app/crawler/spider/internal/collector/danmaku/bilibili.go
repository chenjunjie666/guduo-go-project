package danmaku

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/lib/bilibili/danmaku"
	"guduo/app/crawler/spider/internal/storage"
	"regexp"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func bilibiliHandle() {
	detailUrls := storage.Bilibili.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go bilibiliDanmakuContent(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取哔哩哔哩弹幕内容
func bilibiliDanmakuContent(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()
	unqionWorkerId := fmt.Sprintf("%s_%d", "bilibili_", time.Now().UnixNano())
	oids := common.Bilibili.ParseCids(u)

	danmakuSlice := make([]string, 0, 10000)

	wg2 := &sync.WaitGroup{}
	if len(oids) != 0 {
		for _, oid := range oids {
			wg2.Add(1)
			go func(oid string) {
				historyDanmaku := getHistoryDanmaku(oid, unqionWorkerId)
				danmakuSlice = append(danmakuSlice, historyDanmaku...)
				wg2.Done()
			}(oid)
		}

		wg2.Wait()
		log.Info("bilibili 单剧弹幕采集完毕，开始存储")
		// 存储获取到的简介
		storage.Bilibili.StoreDanmakuContent(danmakuSlice, JobAt, sid)
	} else {
		log.Warn(fmt.Sprintf("链接:%s，未找到视频弹幕", u))
		// 不要return，后续会对0值做处理
		//return
	}

	core.ReleaseVip(unqionWorkerId, common.Bilibili.PlatformId)

	num := bilibiliDanmakuCount(u)
	log.Info("++++++++++++++++++++++++bilibili视频弹幕，showid", sid, " 抓取完毕，共抓取", len(oids), "集，总共", num, "条弹幕，准备开始存储++++++++++++++++++++++++")
	storage.Bilibili.StoreDanmakuCount(num, JobAt, sid)
}

// 获取历史弹幕
func getHistoryDanmaku(oid, unqionWorkerId string) []string {
	c := common.Bilibili.Collector(ModName)

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	findFlag := false

	apiHost := common.Bilibili.ApiHosts.ApiHost
	date := common.Bilibili.ParseDate()
	apiUrl := fmt.Sprintf("%s/x/v2/dm/web/history/seg.so?type=1&oid=%s&date=%s",
		apiHost,
		oid,
		date,
	)

	danmakuSlice := make([]string, 0, 10000)
	mapLock := &sync.Mutex{}
	c.OnResponse(func(r *colly.Response) {
		contentType := r.Headers.Get("Content-Type")
		findFlag = true

		danmakuByte := r.Body

		isOctetStream, _ := regexp.MatchString("octet-stream", contentType)
		isJson, _ := regexp.MatchString("json", contentType)
		if isOctetStream {
			// 弹幕池
			dmPool := danmaku.DmSegMobileReply{}
			// 解析 bytebuffer 到弹幕池中
			err := proto.Unmarshal(danmakuByte, &dmPool)
			if err != nil {
				log.Warn("已找到弹幕原始文件，但解析失败！")
			}

			// 关键字段：
			// row.Content 正文
			// row.Ctime 弹幕发送时间
			for _, row := range dmPool.Elems {
				if filterLastDanmaku(int64(JobAt), row.Ctime) {
					mapLock.Lock()
					ct := row.Content
					danmakuSlice = append(danmakuSlice, ct)
					mapLock.Unlock()
				}
			}
		} else if isJson {
			ErrCode, _ := jsonparser.GetInt(danmakuByte, "code")
			if ErrCode == -101 {
				log.Warn("账号未登录,请检查Cookie!")
			}
		} else {
			log.Warn(fmt.Sprintf("链接:%s有误，请检查！", apiUrl))
		}

	})

	_ = c.Visit(apiUrl, unqionWorkerId)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕文件", apiUrl))
	}
	return danmakuSlice
}

func filterLastDanmaku(jobAt int64, timestamp int64) bool {
	IsLastDanmaku := false
	now := time.Unix(jobAt, 0).Unix()
	startTimeStamp := now - 3600*4
	endTimeStamp := now
	//fmt.Println(startTimeStamp, timestamp, endTimeStamp)
	if timestamp >= startTimeStamp && timestamp <= endTimeStamp {
		IsLastDanmaku = true
	}
	return IsLastDanmaku
}

func bilibiliDanmakuCount(u string) int64 {
	c := common.Bilibili.Collector("弹幕数")

	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})
	num := int64(0)
	findFlag := false
	seasonId := common.Bilibili.ParseSeasonId(u)
	c.OnResponse(func(r *colly.Response) {
		findFlag = true
		dcTemp := r.Body
		num, _ = jsonparser.GetInt(dcTemp, "result", "danmakus")
	})
	url := fmt.Sprintf("%s/pgc/web/season/stat?season_id=%s",
		common.Bilibili.ApiHosts.ApiHost,
		seasonId,
	)
	_ = c.Visit(url)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕数", u))
	}

	//c := common.Bilibili.Collector("弹幕数")
	//num := int64(0)
	//c.OnResponse(func(r *colly.Response) {
	//	reg := regexp.MustCompile(`__INITIAL_STATE__.*?\(function`)
	//	b := reg.Find(r.Body)
	//	s := strings.Trim(string(b), "__INITIAL_STATE__(function ")
	//	num, _ = jsonparser.GetInt([]byte(s), "mediaInfo", "stat", "danmakus")
	//})
	//
	//c.Visit(u)
	return num
}

// 获取当前弹幕池
//func getCurrentDanmakuPool(cid string) []map[string]string {
//	c := common.Bilibili.Collector(ModName)
//
//	c.OnError(func(r *colly.Response, e error) {
//		c.Retry(r, e)
//	})
//
//	findFlag := false
//
//	var danmakuSlice []map[string]string
//
//	apiHost := common.Bilibili.ApiHosts.ApiHost
//	// 经观察此参数为固定参数值
//	segmentIndex := 1
//	apiUrl := fmt.Sprintf("%s/x/v2/dm/web/seg.so?type=1&oid=%s&segment_index=%d",
//		apiHost,
//		cid,
//		segmentIndex,
//	)
//
//	c.OnResponse(func(r *colly.Response) {
//		contentType := r.Headers.Get("Content-Type")
//		findFlag = true
//
//		danmakuByte := r.Body
//
//		isOctetStream, _ := regexp.MatchString("octet-stream", contentType)
//		isJson, _ := regexp.MatchString("json", contentType)
//		if isOctetStream {
//			// 弹幕池
//			dmPool := danmaku.DmSegMobileReply{}
//			// 解析 bytebuffer 到弹幕池中
//			err := proto.Unmarshal(danmakuByte, &dmPool)
//			if err != nil {
//				log.Warn(fmt.Println("已找到弹幕原始文件，但解析失败！"))
//			}
//
//			// 关键字段：
//			// row.Content 正文
//			// row.Ctime 弹幕发送时间
//			for _, row := range dmPool.Elems {
//				var danmakuMap map[string]string
//				danmakuMap = make(map[string]string)
//				if filterLastDanmaku(int64(JobAt), row.Ctime) {
//					danmakuMap["Content"] = row.Content
//					danmakuMap["Ctime"] = strconv.FormatInt(row.Ctime, 10)
//					danmakuSlice = append(danmakuSlice, danmakuMap)
//				}
//			}
//
//			fmt.Println(danmakuSlice)
//
//		} else if isJson {
//			ErrCode, _ := jsonparser.GetInt(danmakuByte, "code")
//			if ErrCode == -101 {
//				log.Warn(fmt.Println("账号未登录,请检查Cookie!"))
//			}
//		} else {
//			log.Warn(fmt.Sprintf("链接:%s有误，请检查！", apiUrl))
//		}
//
//	})
//
//	_ = c.Visit(apiUrl)
//
//	if !findFlag {
//		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕文件", apiUrl))
//	}
//	return danmakuSlice
//}
