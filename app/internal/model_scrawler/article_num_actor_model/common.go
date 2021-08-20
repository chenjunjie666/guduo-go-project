package article_num_actor_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean/article_count_actor_daily_model"
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
func SaveArticleNum(an int64, ja uint, aid, pid uint64)  {
	// 同一个周期内，一个剧集，只应该抓取一次
	r := Model().Where("actor_id = ? and job_at = ?", aid, ja).Limit(1).Find(nil)

	if r.RowsAffected > 0 {
		return
	}

	d := &Table{
		PlatformId: pid,
		ActorId: aid,
		Num: an,
		JobAt: ja,
	}

	Model().Create(&d)

	article_count_actor_daily_model.SaveCurCount(an, ja, aid, pid)
}