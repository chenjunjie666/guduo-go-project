package hot_model

import (
	"gorm.io/gorm"
	hot_daily_model "guduo/app/internal/model_clean/hot_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存热度
func SaveHotCount(hot int64, jobAt uint, showId, platformId uint64) {
	var cnt int64
	Model().Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Count(&cnt)

	if cnt > 0{
		return
	}

	d := Table{
		ShowId: showId,
		PlatformId: platformId,
		Hot: hot,
		JobAt: jobAt,
	}

	Model().Create(&d)
	hot_daily_model.SaveCurHot(hot, jobAt, showId, platformId)
}