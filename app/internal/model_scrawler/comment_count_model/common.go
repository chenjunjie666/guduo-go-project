package comment_count_model

import (
	"gorm.io/gorm"
	comment_count_daily_model "guduo/app/internal/model_clean/comment_count_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存演员
func SaveCommentCount(num int64, jobAt uint, showId, platformId uint64) {
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

	comment_count_daily_model.SaveCurCount(num, jobAt, showId, platformId)
}