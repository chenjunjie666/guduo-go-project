package indicator_actor_model

import (
	"gorm.io/gorm"
	indicator_actor_daily_model "guduo/app/internal/model_clean/indicator_actor_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存总指数
func SaveIndicator(num int64, jobAt uint, ActorId, platformId uint64) {
	var cnt int64
	Model().Where("job_at = ? and actor_id = ? and platform_id = ?", jobAt, ActorId, platformId).
		Count(&cnt)

	if cnt > 0{
		return
	}

	d := Table{
		ActorId: ActorId,
		PlatformId: platformId,
		Num: num,
		JobAt: jobAt,
	}
	Model().Create(&d)
	indicator_actor_daily_model.SaveCurIndicator(num, jobAt, ActorId, platformId)
}