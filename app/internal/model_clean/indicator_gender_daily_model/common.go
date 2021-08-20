package indicator_gender_daily_model

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

// 保存当日当前性别分布
func SaveCurGender(rating map[string]float64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ShowId:       sid,
		PlatformId:   pid,
		MaleRating:   rating["male"],
		FemaleRating: rating["female"],
		DayAt:        da,
	}

	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)
	Model().Create(&row)
}
