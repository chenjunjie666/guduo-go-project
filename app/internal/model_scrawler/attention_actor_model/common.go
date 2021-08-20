package attention_actor_model

import (
	"gorm.io/gorm"
	attention_actor_daily_model "guduo/app/internal/model_clean/attention_actor_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}


func StoreAttention(a, p int64, ja uint, aid, pid uint64){
	// 判断重复
	var cnt int64
	Model().Where("job_at = ?", ja).
		Where("actor_id = ? and platform_id = ?", aid, pid).
		Count(&cnt)
	if cnt > 0 {
		return
	}

	d := &Table{
		ActorId:     aid,
		PlatformId: pid,
		Attention:  a,
		Post:       p,
		JobAt:      ja,
	}

	Model().Create(&d)

	attention_actor_daily_model.SaveCurAttention(a, p, ja, aid, pid)
}