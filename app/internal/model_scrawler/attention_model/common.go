package attention_model

import (
	"gorm.io/gorm"
	attention_daily_model "guduo/app/internal/model_clean/attention_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}


func StoreAttention(a, p int64, ja uint, sid, pid uint64){
	// 判断重复
	var cnt int64
	Model().Where("job_at = ?", ja).
		Where("show_id = ? and platform_id = ?", sid, pid).
		Count(&cnt)
	if cnt > 0 {
		return
	}

	d := &Table{
		ShowId:     sid,
		PlatformId: pid,
		Attention:  a,
		Post:       p,
		JobAt:      ja,
	}

	Model().Create(&d)

	attention_daily_model.SaveCurAttention(a, p, ja, sid, pid)
}