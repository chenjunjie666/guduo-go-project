package short_comment_model

import (
	"gorm.io/gorm"
	short_comment_count_daily_model "guduo/app/internal/model_clean/short_comment_count_daily_model"
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
func SaveShortCommentCount(num int64, jobAt uint, showId, platformId uint64) {
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

	short_comment_count_daily_model.SaveCurShortCommentCount(num, jobAt, showId, platformId)
}