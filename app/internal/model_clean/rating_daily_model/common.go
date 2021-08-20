package rating_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前评分
func SaveCurPlayCount(num float64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Rating: num,
		DayAt: da,
	}

	r := Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Limit(1).Find(&row)
	if r.RowsAffected > 0 {
		r.Updates(row)
	}else{
		Model().Create(&row)
	}
}