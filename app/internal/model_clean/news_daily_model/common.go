package news_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/app/internal/model_clean/news_trend_daily_model"
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

// 保存当日当前新闻数
func SaveCurNews(num int64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

	lastRow := &Table{}
	Model().Select("num").Where("show_id = ? and platform_id = ?", sid, pid).Order("day_at desc").Limit(1).Find(&lastRow)

	if num == 0 {
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
	news_trend_daily_model.SaveCurNews(trend, da, sid, pid)
}


// 获取单日的文章数
func GetNewsCount(sid uint64, day []uint, pid... uint64) int64 {
	mm := Model()
	mm = mm.Select("sum(IF(custom_num != 0, custom_num, num)) as total_count")
	mm = mm.Where("show_id = ?", sid)

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

	r := mm.Find(&res)

	if r.RowsAffected > 0 {
		return res.TotalCount
	}

	return 0
}