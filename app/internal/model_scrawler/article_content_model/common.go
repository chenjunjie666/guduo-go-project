package article_content_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/constant"
	article_content_current_model "guduo/app/internal/model_clean/article_content_current_model"
	"guduo/pkg/db"
	"guduo/pkg/time"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}


func SaveContent(uid, name, _time, content string, forward int64, ja uint, sid, pid uint64) int64 {
	// 先检测这个热门文章的唯一id是否存在
	var cnt int64
	Model().Where("uid = ? and platform_id = ?", uid, constant.PlatformIdWeibo).Count(&cnt)
	if cnt > 0 {
		return 1
	}

	publishAt := time.YYYYmmddToSecTimestamp(_time)
	d := &Table{
		PlatformId: pid,
		ShowId: sid,
		UID: uid,
		Content: content,
		Author: name,
		PublishAt: publishAt,
		Forward: forward,
		JobAt: ja,
	}

	Model().Create(&d)

	article_content_current_model.SaveCurContent(name, content, forward, publishAt, ja, sid, pid)
	return 0
}