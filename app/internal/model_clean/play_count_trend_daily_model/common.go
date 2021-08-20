package play_count_trend_daily_model

import (
	"gorm.io/gorm"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前播放量
func SaveCurPlayCount(num int64, da uint, sid, pid uint64) {
	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: num,
		DayAt: da,
	}

	Model().Create(&row)
}