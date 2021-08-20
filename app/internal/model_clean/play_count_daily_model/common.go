package play_count_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/app/internal/model_clean/play_count_trend_daily_model"
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
func SaveCurPlayCount(num int64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

	lastRow := &Table{}
	Model().Select("num").Where("show_id = ? and platform_id = ?", sid, pid).Order("day_at desc").Limit(1).Find(&lastRow)

	if num < lastRow.Num {
		num = lastRow.Num
	}

	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: num,
		DayAt: da,
	}
	Model().Create(&row)

	trend := num - lastRow.Num
	play_count_trend_daily_model.SaveCurPlayCount(trend, da, sid, pid)
}