package indicator_daily_model

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

// 保存当日当前主要指数
func SaveCurIndicator(i float64, ja uint, sid, pid uint64) {
	// 0 理论不应该，所以0就直接不存
	if i == 0 {
		return
	}

	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: i,
		DayAt: da,
	}

	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)
	Model().Create(&row)
}


// 获取单日的平台指数
func GetIndicator(sid []uint64, day []uint, pid... uint64) []*Table {
	mm := Model()
	mm = mm.Select("SUM(IF(custom_num != 0, custom_num, num)) as num", "show_id")
	mm = mm.Where("show_id IN ?", sid)

	if len(day) == 1 {
		mm = mm.Where("day_at = ?", day)
	}else if len(day) == 2 {
		mm = mm.Where("day_at BETWEEN ? AND ?", day[0], day[1] -1 /* 取值范围不包含结束时间，所以这里-1秒 */)
	}

	if len(pid) > 0 {
		mm = mm.Where("platform_id IN ?", pid)
	}

	var res []*Table

	mm.Group("show_id").Find(&res)

	return res
}