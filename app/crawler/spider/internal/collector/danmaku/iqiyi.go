package danmaku

import (
	"fmt"
	"guduo/app/crawler/spider/internal/collector/common"
	"guduo/app/crawler/spider/internal/lib/iqiyi/danmaku"
	"guduo/app/crawler/spider/internal/storage"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

func iqiyiHandle() {
	detailUrls := storage.Iqiyi.GetDetailUrl()

	for _, row := range detailUrls {
		wg.Add(1)
		ch.PushJob()
		go iqiyiDanmakuContent(row.Url, row.ShowId)
	}

	wg.Done()
}

// 从页面解析并获取获取爱奇艺弹幕数
func iqiyiDanmakuContent(u string, sid uint64) {
	defer ch.PopJob()
	defer wg.Done()
	c := common.Iqiyi.Collector("爱奇艺" + ModName)
	c.Async = true

	unqionWorkerId := fmt.Sprintf("%s_%d", "iqiyi_", time.Now().UnixNano())
	vids := common.Iqiyi.ParseVids(u)
	if len(vids) == 0 {
		log.Warn(fmt.Sprintf("没有解析到vids，可能是院线电影，无弹幕，链接:%s", u))
		return
	}

	tc := int64(0) // 总弹幕数

	c.OnError(func(r *colly.Response, e error) {
		if r.StatusCode == 404 {
			return
		}
		c.Retry(r, e)
	})

	findFlag := false

	danmakuSlice := make([]string, 0, 5000)
	mapLock := &sync.Mutex{}
	c.OnResponse(func(r *colly.Response) {
		findFlag = true
		dcTemp := r.Body
		resList := danmaku.Decode(dcTemp)
		for _, res := range resList {
			tc++
			//if filterLastDanmaku(int64(JobAt), res.Ctime/1e9) {
				mapLock.Lock()
				ct := res.Content
				danmakuSlice = append(danmakuSlice, ct)
				mapLock.Unlock()
			//}
		}
	})

	lenTmp := common.Iqiyi.ParseLength(u)
	// 爱奇艺弹幕的url的时间step是300秒
	for _, vid := range vids {
		length := int(lenTmp[vid])
		for i := 1; i <= length; i++ {
			url := fmt.Sprintf("%s/bullet/%s/%s/%s_300_%d.z",
				common.Iqiyi.ApiHosts.BulletHost,
				vid[len(vid)-4:len(vid)-2],
				vid[len(vid)-2:],
				vid,
				i,
			)
			_ = c.Visit(url, unqionWorkerId)
		}
	}

	c.Wait()

	c.ReleaseVip(unqionWorkerId)
	log.Info("++++++++++++++++++++++++爱奇艺视频弹幕，showid", sid, " 抓取完毕，共抓取", len(vids), "集，总共", tc, "条弹幕，准备开始存储++++++++++++++++++++++++")

	if !findFlag {
		log.Warn(fmt.Sprintf("链接:%s，未找到弹幕数", u))
		//return
	}

	storage.Iqiyi.StoreDanmakuContent(danmakuSlice, JobAt, sid)
	storage.Iqiyi.StoreDanmakuCount(tc, JobAt, sid)
}
