package half_month

import (
	time2 "guduo/pkg/time"
	"sync"
	"time"
)
const ModName = "每周指标"

var JobAt = time2.Today()
var StartAt uint
var EndAt uint
var Day uint
var lastMonth uint

var wg = &sync.WaitGroup{}

func Run() {
	t := time.Now()
	firstday := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)

	StartAt = uint(firstday.Unix())
	EndAt = time2.Today()
	lastMonth = uint(firstday.AddDate(0, -1, 0).Unix())

	// 15号计算，14号不计算，所以14 - 1 = 13
	if JobAt - StartAt <= (86400 * 13) {
		return
	}
	Day = JobAt - StartAt + 1


	//jobNum := 2
	//wg.Add(jobNum)

	// 骨朵剧集热度
	guduoHotHandle()
	guduoActorHotHandle()
	guduoActorDomiHandle()
	//wg.Wait()

	// 每月电影播放量
	moviePlayCountHandle()
}