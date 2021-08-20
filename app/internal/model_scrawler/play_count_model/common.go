package play_count_model

import (
	"gorm.io/gorm"
	play_count_daily_model "guduo/app/internal/model_clean/play_count_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}


func StorePlayCount(pc int64, ja uint, sid, pid uint64){
	var cnt int64
	Model().Where("job_at = ? and show_id = ? and platform_id = ?", ja, sid, pid).
		Count(&cnt)

	if cnt > 0{
		return
	}

	d := Table{
		ShowId: sid,
		PlatformId: pid,
		Num: pc,
		JobAt: ja,
	}
	Model().Create(&d)

	play_count_daily_model.SaveCurPlayCount(pc, ja, sid, pid)
}