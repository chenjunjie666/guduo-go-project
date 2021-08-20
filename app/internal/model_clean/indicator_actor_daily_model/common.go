package indicator_actor_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/pkg/db"
	"guduo/pkg/model"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前主要指数
func SaveCurIndicator(i int64, ja uint, aid, pid uint64) {
	// 0 理论不应该，所以0就直接不存
	if i == 0 {
		return
	}

	da := helper.JobAt2DayAt(ja)
	row := &Table{
		ActorId: aid,
		PlatformId: pid,
		Num: i,
		DayAt: da,
	}
	Model().Where("actor_id = ? and platform_id = ? and day_at = ?", aid, pid, da).Delete(nil)
	Model().Create(&row)
}


// 获取单日的平台指数
func GetIndicator(aid uint64, day []uint, pid... uint64) int64 {
	mm := Model()
	mm = mm.Select("sum(IF(custom_num != 0, custom_num, num)) as total_count")
	mm = mm.Where("actor_id = ?", aid)

	if len(day) == 1 {
		mm = mm.Where("day_at = ?", day)
	}else if len(day) == 2 {
		mm = mm.Where("day_at BETWEEN ? AND ?", day[0], day[1] -1 /* 取值范围不包含结束时间，所以这里-1秒 */)
	}

	if len(pid) > 0 {
		mm = mm.Where("platform_id IN ?", pid)
	}

	res := &struct{
		TotalCount model.Int
	}{}

	r := mm.Find(res)

	if r.RowsAffected > 0 {
		return res.TotalCount
	}

	return 0
}