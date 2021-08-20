package article_count_actor_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/app/internal/model_clean/article_count_actor_trend_daily_model"
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
func SaveCurCount(cc int64, ja uint, aid, pid uint64) {
	da := helper.JobAt2DayAt(ja)

	Model().Where("actor_id = ? and platform_id = ? and day_at = ?", aid, pid, da).Delete(nil)

	lastRow := &Table{}
	Model().Select("num").Where("actor_id = ? and platform_id = ?", aid, pid).Order("day_at desc").Limit(1).Find(&lastRow)

	row := &Table{
		ActorId: aid,
		PlatformId: pid,
		Num: cc,
		DayAt: da,
	}
	Model().Create(&row)

	trend := cc - lastRow.Num
	article_count_actor_trend_daily_model.SaveCurCount(trend, da, aid, pid)
}


// 获取文章数
func GetArticleNum(aid uint64, day []uint, pid... uint64) int64 {
	mm := Model()
	mm = mm.Where("actor_id = ?", aid)

	if len(day) == 1 {
		mm = mm.Select("IF(custom_num != 0, custom_num, num) as total_count")
		mm = mm.Where("day_at = ?", day)
	}else if len(day) == 2 {
		mm = mm.Select("sum(IF(custom_num != 0, custom_num, num)) as total_count")
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