package half_week

import (
	time2 "guduo/pkg/time"
	"time"
)
const ModName = "每周指标"

var JobAt = time2.Today()
var StartAt uint
var EndAt uint
var Day uint
var lastWeek uint

//var wg = &sync.WaitGroup{}

func Run() {
	week := uint(time.Now().Weekday())

	// 周一周二不跑周指标
	if week == 1 || week == 2 {
		return
	}

	if week == 0 {
		week = 7
	}

	Day = week
	EndAt = time2.Today()
	StartAt = EndAt - (week - 1) * 86400


	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastWeek = uint(weekStart.AddDate(0, 0, -7).Unix())

	//jobNum := 2
	//wg.Add(jobNum)

	// 骨朵剧集热度
	guduoHotHandle()
	guduoActorHotHandle()
	guduoActorDomiHandle()
	//wg.Wait()

	// 每周电影播放量
	moviePlayCountHandle()
}