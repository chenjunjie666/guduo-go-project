package comment_count_trend_daily_model

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

// 保存当日当前评论数
func SaveCurCount(cc int64, da uint, sid, pid uint64) {
	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: cc,
		DayAt: da,
	}

	Model().Create(&row)
}

// 获取单日的评论总数
func GetCommentCount(sid []uint64, day []uint, pid... uint64) []*Table {
	mm := Model()
	mm = mm.Select("sum(IF(custom_num != 0, custom_num, num)) as num","show_id")
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