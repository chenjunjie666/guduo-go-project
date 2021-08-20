package news_model

import (
	"gorm.io/gorm"
	news_daily_model "guduo/app/internal/model_clean/news_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存粉丝数
func SaveNewsNum(num int64, jobAt uint, showId, platformId uint64) {
	var cnt int64
	Model().Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Count(&cnt)

	if cnt > 0{
		return
	}

	d := Table{
		ShowId: showId,
		PlatformId: platformId,
		Num: num,
		JobAt: jobAt,
	}
	Model().Create(&d)

	news_daily_model.SaveCurNews(num, jobAt, showId, platformId)
}