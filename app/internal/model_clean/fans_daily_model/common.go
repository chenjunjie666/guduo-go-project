package fans_daily_model

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

// 保存当日当前粉丝数
func SaveCurFans(fs int64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: fs,
		DayAt: da,
	}

	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)
	Model().Create(&row)
}