package article_num_model

import (
	"gorm.io/gorm"
	article_num_daily_model "guduo/app/internal/model_clean/article_count_daily_model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存文章数
func SaveArticleNum(an int64, ja uint, sid, pid uint64)  {
	// 同一个周期内，一个剧集，只应该抓取一次
	var cnt int64
	Model().Where("show_id = ? and job_at = ?", sid, ja).Count(&cnt)

	if cnt > 0 {
		return
	}

	d := &Table{
		PlatformId: pid,
		ShowId: sid,
		Num: an,
		JobAt: ja,
	}

	Model().Create(&d)

	article_num_daily_model.SaveCurCount(an, ja, sid, pid)
}