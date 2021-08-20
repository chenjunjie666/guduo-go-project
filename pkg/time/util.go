package time

import time2 "time"

// 获取今天的0点的时间戳
func Today() uint {
	cur := time2.Now()
	day := time2.Date(cur.Year(), cur.Month(), cur.Day(), 0, 0, 0, 0, time2.Local).Unix()

	return uint(day)
}
