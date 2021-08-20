package helper

import (
	time2 "guduo/pkg/time"
	"strconv"
	"time"
)

// 获取周期-整点运行
func GetJobAt() uint {
	return time2.YYYYmmddToSecTimestamp(time.Now().Format("2006-01-02 15") + ":00:00") // 格式化为秒级时间为0的时间戳
}

// 获取周期-有半小时单位的
func GetJobAtHalf() uint {
	tNow := time.Now()

	t1 := tNow.Format("2006-01-02 15")
	min, _ := strconv.ParseInt(tNow.Format("04"), 10, 64)

	t2 := ""
	if min >= 30 {
		t2 = ":30:00"
	}else{
		t2 = ":00:00"
	}

	t := t1 + t2
	return time2.YYYYmmddToSecTimestamp(t)
}



func JobAt2DayAt(ja uint) uint {
	dayStr := time.Unix(int64(ja), 0).Format("2006-01-02") + " 00:00:00"
	dayAt := time2.YYYYmmddToSecTimestamp(dayStr)

	return dayAt
}