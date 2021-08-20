package rating_model

import (
	"gorm.io/gorm"
	rating_daily_model "guduo/app/internal/model_clean/rating_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当前评分
func SaveRatingNum(num float64, jobAt uint, showId, platformId uint64) {
	r := Model().Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Limit(1).Find(nil)

	if r.RowsAffected > 0{
		return
	}

	d := Table{
		ShowId: showId,
		PlatformId: platformId,
		Num: num,
		JobAt: jobAt,
	}
	Model().Create(&d)

	rating_daily_model.SaveCurPlayCount(num, jobAt, showId, platformId)
}