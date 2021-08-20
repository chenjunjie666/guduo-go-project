package danmaku_count_model

import (
	"gorm.io/gorm"
	danmaku_count_daily_model "guduo/app/internal/model_clean/danmaku_count_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存弹幕数
func SaveDanmakuCount(num int64, jobAt uint, showId, platformId uint64) {
	Model().Select("num").Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Delete(nil)

	d := Table{
		ShowId: showId,
		PlatformId: platformId,
		Num: num,
		JobAt: jobAt,
	}
	Model().Create(&d)

	danmaku_count_daily_model.SaveCurCount(num, jobAt, showId, platformId)
}