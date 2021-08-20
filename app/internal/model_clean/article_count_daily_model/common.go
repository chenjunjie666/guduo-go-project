package article_count_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/constant"
	"guduo/app/internal/helper"
	"guduo/app/internal/model_clean/article_count_trend_daily_model"
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

// 保存当日当前文章数
func SaveCurCount(cc int64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)


	lastRow := &Table{}

	trend := int64(0)
	if pid == constant.PlatformIdWeibo {
		trend = cc
		// 微博由于每次爬取的都是一个增量，所以需要先取上次数据，相加得到一个总数，然后再删除当日记录
		Model().Select("num").Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Limit(1).Find(&lastRow)
		Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)
		// 微博不是爬取的总数，他总是爬取变化量
		cc += lastRow.Num
	}else {
		// 其他平台抓的总数所以先删掉当日总数，取上次的总数，然后相减得出变化量
		Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

		Model().Select("num").Where("show_id = ? and platform_id = ?", sid, pid).Order("day_at desc").Limit(1).Find(&lastRow)

		trend = cc - lastRow.Num
	}

	row := &Table{
		ShowId: sid,
		PlatformId: pid,
		Num: cc,
		DayAt: da,
	}
	Model().Create(&row)

	article_count_trend_daily_model.SaveCurCount(trend, da, sid, pid)
}


// 获取文章数
func GetArticleNum(sid uint64, day []uint, pid... uint64) int64 {
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