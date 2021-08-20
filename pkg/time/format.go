package time

import (
	"regexp"
	"time"
)

// 精确到秒
const LayoutYmdHis = "2006-01-02 15:04:05"
// 精确到天
const LayoutYmd = "2006-01-02"
// 月日
const LayoutMd = "01-02"

// 把 Y-m-d H:i:s 格式转为秒级时间戳
func YYYYmmddToSecTimestamp(ts string) uint {
	reg := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	b := reg.MatchString(ts)


	if !b {
		return 0
	}

	tl, _ := time.LoadLocation("Local")
	t, e := time.ParseInLocation(LayoutYmdHis, ts, tl)
	if e != nil {
		return 0
	}

	return uint(t.Unix())
}

func Time() uint {
	return uint(time.Now().Unix())
}


func TimeToStr(layout string, time_ uint) string {
	t := int64(time_)
	return time.Unix(t, 0).Format(layout)
}