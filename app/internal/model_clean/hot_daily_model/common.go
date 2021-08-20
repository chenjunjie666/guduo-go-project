package hot_daily_model

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

// 保存当日当前视频热度
func SaveCurHot(hot int64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: hot,
		DayAt: da,
	}

	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)
	var lastRow *Table
	Model().Select("num").Where("show_id = ? and platform_id = ?", sid, pid).Order("day_at desc").Limit(1).Find(&lastRow)

	if hot == 0 {
		hot = lastRow.Num
	}

	Model().Create(&row)
}