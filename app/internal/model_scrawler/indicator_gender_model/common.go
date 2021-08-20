package indicator_gender_model

import (
	indicator_gender_daily_model "guduo/app/internal/model_clean/indicator_gender_daily_model"
	"guduo/pkg/db"

	"gorm.io/gorm"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存总指数
func SaveIndicatorGender(rating map[string]float64, jobAt uint, showId, platformId uint64) {
	var cnt int64
	Model().Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Count(&cnt)

	if cnt > 0 {
		return
	}

	d := Table{
		ShowId:       showId,
		PlatformId:   platformId,
		MaleRating:   rating["male"],
		FemaleRating: rating["female"],
		JobAt:        jobAt,
	}
	Model().Create(&d)

	indicator_gender_daily_model.SaveCurGender(rating, jobAt, showId, platformId)
}
