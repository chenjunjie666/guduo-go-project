package comment_count_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/app/internal/model_clean/comment_count_trend_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前评论数
func SaveCurCount(cc int64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

	lastRow := &Table{}
	Model().Select("num").Where("show_id = ? and platform_id = ?", sid, pid).Order("day_at desc").Limit(1).Find(&lastRow)

	if cc < lastRow.Num {
		cc = lastRow.Num
	}

	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: cc,
		DayAt: da,
	}
	Model().Create(&row)

	trend := cc - lastRow.Num
	comment_count_trend_daily_model.SaveCurCount(trend, da, sid, pid)
}