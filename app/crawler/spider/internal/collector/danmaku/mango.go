package danmaku

import (
	"bytes"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/storage"
	"strings"
	"sync"
	"time"
)

func mangoHandle() {

	detailUrls := storage.Mango.GetDetailUrl()

	for _, row := range detailUrls {
		// 不满足条件的url，直接跳过，不要再去浪费资源
		if !strings.Contains(row.Url, "com/b/") {
			continue
		}
		wg.Add(1)
		ch.PushJob()
		go mangoDanmakuContent(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取芒果TV弹幕内容
func mangoDanmakuContent(u string, sid uint64) {
	defer wg.Done()
	unqionWorkerId := fmt.Sprintf("%s_%d", "mgtv_", time.Now().UnixNano())

	findFlag := false

	vids := common.Mango.ParseVids(u)
	cid := common.Mango.ParseCid(u)

	lengthAll := common.Mango.ParseLen(u)


	//q, _ := queue.New(
	//	1, // Number of consumer threads
	//	&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	//)

	tc := int64(0)
	tidx := 0

	c := common.Mango.Collector("芒果" + ModName)

	c.Async = true

	c.OnError(func(r *colly.Response, e error) {
		if r.StatusCode == 404 {
			return
		}
		c.Retry(r, e)
	})

	danmakuSlice := make([]string, 0, 10000)
	danmakuIdSlice := make([]string, 0, 10000)
	mapLock := &sync.Mutex{}
	c.OnResponse(func(r *colly.Response) {
		tidx++
		findFlag = true
		dcTemp := r.Body
		dcTemp = bytes.Trim(dcTemp, "jsonp_1234567890();")

		//msg, _ := jsonparser.GetString(dcTemp, "msg")
		rowTc := int64(0)
		jsonparser.ArrayEach(dcTemp, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			rowTc++
			mapLock.Lock()
			id, _ := jsonparser.GetString(value, "ids")
			ct, _ := jsonparser.GetString(value, "content")
			danmakuSlice = append(danmakuSlice, ct)
			danmakuIdSlice = append(danmakuIdSlice, id)
			mapLock.Unlock()

		}, "data", "items")

		tc += rowTc
	})

	//https://galaxy.bz.mgtv.com/cdn/opbarrage?vid=%s&cid=%s&time=%d
	//http://galaxy.person.mgtv.com/rdbarrage?vid=%&cid=%s&time=%d
	wg2 := &sync.WaitGroup{}
	urlBaseMap := make(map[string]string)
	for _, vid := range vids {
		wg2.Add(1)
		go getMangoDanmakuCDNUrl(vid, cid, unqionWorkerId, urlBaseMap, wg2)
	}

	wg2.Wait()

	for _, vid := range vids {
		length := lengthAll[vid]
		interval := 60000
		apiUrlBase := urlBaseMap[vid]
		if apiUrlBase == "" {
			for i := 0; i <= int(length/60); i++ {
				url := fmt.Sprintf("%s/cdn/opbarrage?vid=%s&cid=%s&time=%d",
					common.Mango.ApiHosts.GalaxyBxHost,
					vid,
					cid,
					i*interval,
				)
				c.Visit(url, unqionWorkerId)
			}
		} else {
			for i := 0; i <= int(length/60); i++ {
				url := fmt.Sprintf("%s/%d.json",
					apiUrlBase,
					i,
				)
				c.Visit(url, unqionWorkerId)
			}
		}
	}
	c.Wait()

	log.Info("++++++++++++++++++++++++芒果TV视频弹幕，showid", sid, " 抓取完毕，共抓取", len(vids), "集，总共", tc, "条弹幕，准备开始存储++++++++++++++++++++++++")
	// 释放通道
	c.ReleaseVip(unqionWorkerId)

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕数，show_id:%d", u, sid))
		//return
	}

	// 芒果tv不要等到数据库操作完了再pop，太慢了
	ch.PopJob()

	storage.Mango.StoreDanmakuContent(danmakuSlice, JobAt, sid, danmakuIdSlice)
	storage.Mango.StoreDanmakuCount(tc, JobAt, sid)
}

var MangoMapLock = &sync.Mutex{}

func getMangoDanmakuCDNUrl(vid, cid, unqionWorkerId string, urlBaseMap map[string]string, wg2 *sync.WaitGroup) {
	defer wg2.Done()
	version := "3.0.0"
	deviceid := "b02aa535-c3de-43b7-b705-99f4fd793a7c"
	appVersion := "3.0.0"
	jp := common.Mango.GenJsonpCallbackStr()

	danmakuUrl := ""
	apiUrl := fmt.Sprintf("https://galaxy.bz.mgtv.com/getctlbarrage?version=%s&vid=%s&abroad=0&pid=0&os=&uuid=&deviceid=%s&cid=%s&ticket=&mac=&platform=0&appVersion=%s&reqtype=form-post&callback=%s",
		version,
		vid,
		deviceid,
		cid,
		appVersion,
		jp,
	)

	c := common.Mango.Collector("弹幕CDN")
	c.OnError(func(r *colly.Response, e error) {
		c.Retry(r, e)
	})

	c.OnResponse(func(r *colly.Response) {
		json_ := bytes.Trim(r.Body, jp+"();")

		link, _ := jsonparser.GetString(json_, "data", "cdn_version")

		if link == "" {
			return
		}

		danmakuUrl = fmt.Sprintf("https://bullet-ws.hitv.com/%s", link)
	})

	c.Visit(apiUrl, unqionWorkerId)
	MangoMapLock.Lock()
	urlBaseMap[vid] = danmakuUrl
	MangoMapLock.Unlock()
}
