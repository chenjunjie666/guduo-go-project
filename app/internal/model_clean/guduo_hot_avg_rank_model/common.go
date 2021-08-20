package guduo_hot_avg_rank_model

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
func SaveAvgHot(fs float64, rank int64, ja uint, sid uint64) {
	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ShowId: sid,
		Rank: rank,
		Num: fs,
		DayAt: da,
	}

	Model().Where("show_id", sid).Delete(nil)
	Model().Create(&row)
}