package cfg

import (
	"github.com/gin-gonic/gin"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/weixin_mini/internal/hepler/resp"
	time2 "guduo/pkg/time"
	"strconv"
	"time"
)

func Config(c *gin.Context) {
	platform := constant.GetPlatformMap()
	playType := show_actor_model.GetPlayTypeMap()

	ret := map[string]interface{}{
		"platform": platform,
		"actor_play_type": playType,
	}
	resp.Success(c, ret)
}

func DateConfig(c *gin.Context){
	actor_date := c.Query("actor")
	flag := false
	if actor_date == "1"{
		flag = true
	}

	//fmt.Println(actor_date)
	//panic("--------")


	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	curDay, _ := strconv.Atoi(now.Format("02"))
	curWeek := int(now.Weekday())

	start := time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local)
	endMonth := today
	endWeek := today
	endDay := today


	if curDay <= 15 || flag == true {
		endMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	}
	if (curWeek >= 1 && curWeek <= 3) || flag == true{
		endWeek = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -curWeek)
	}

	if flag == true {
		endDay = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	}

	listYear := make(map[int][]interface{})
	d := start
	idx := 0
	endLoop := false
	for {
		endThisWeek := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 6)

		if endThisWeek.Unix() >= endWeek.Unix() {
			endLoop = true
			endThisWeek = endWeek
		}
		row := []interface{}{
			idx,
			time2.TimeToStr(time2.LayoutMd, uint(d.Unix())),
			time2.TimeToStr(time2.LayoutMd, uint(endThisWeek.Unix())),
		}
		if d.Year() != endThisWeek.Year() {
			row[0] = 0
			idx = 0
			listYear[endThisWeek.Year()] = make([]interface{}, 0, 60)
			listYear[endThisWeek.Year()] = append(listYear[endThisWeek.Year()], row)
		}else{
			listYear[d.Year()] = append(listYear[d.Year()], row)

		}
		if endLoop {
			break
		}
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 7)
		if _, ok := listYear[d.Year()]; !ok {
			listYear[d.Year()] = make([]interface{}, 0, 60)
		}
		idx++
	}


	if flag == true {
		endDay = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	}

	dayPart := map[string]interface{}{
		"start_date": time2.TimeToStr(time2.LayoutYmd, uint(start.Unix())),
		"end_date": time2.TimeToStr(time2.LayoutYmd, uint(endDay.Unix())),
	}

	weekPart := map[string]interface{}{
		"start_date": time2.TimeToStr(time2.LayoutYmd, uint(start.Unix())),
		"end_date": time2.TimeToStr(time2.LayoutYmd, uint(endWeek.Unix())),
		"list": listYear,
	}

	monthPart := map[string]string{
		"start_date": time2.TimeToStr(time2.LayoutYmd, uint(start.Unix())),
		"end_date": time2.TimeToStr(time2.LayoutYmd, uint(endMonth.Unix())),
	}
	yearPart := map[string]string{
		"start_date": time2.TimeToStr(time2.LayoutYmd, uint(start.Unix())),
		"end_date": time2.TimeToStr(time2.LayoutYmd, uint(endDay.Unix())),
	}
	totalPart := map[string]string{
		"start_date": time2.TimeToStr(time2.LayoutYmd, uint(start.Unix())),
		"end_date": time2.TimeToStr(time2.LayoutYmd, uint(endDay.Unix()) - 86400),
	}

	dateCfg := map[string]interface{}{
		"day": dayPart,
		"week": weekPart,
		"month": monthPart,
		"year": yearPart,
		"total": totalPart,
	}

	resp.Success(c, dateCfg)
}

