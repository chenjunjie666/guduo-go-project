package year

import (
	"guduo/app/crawler/clean/internal/core"
	time2 "guduo/pkg/time"
	"time"
)
const ModName = "每年指标"

var JobAt = time2.Today()
var StartAt uint
var EndAt uint
var Day uint

//var wg = &sync.WaitGroup{}

func Run() {
	t := time.Now()
	if t.Day() != 1 || t.Month() != 1 {
		return
	}

	StartAt = uint(time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.Local).AddDate(-1, 0, 0).Unix())
	EndAt = uint(t.Unix() - 86400)

	core.Init()

	//jobNum := 2
	//wg.Add(jobNum)

	// 骨朵剧集热度
	//guduoHotHandle()
	//daily.YearTotalPlayCountHandle()
	//wg.Wait()
}