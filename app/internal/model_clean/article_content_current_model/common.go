package article_content_current_model

import (
	"gorm.io/gorm"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前热门微博
func SaveCurContent(name, content string, forward int64, _time, ja uint, sid, pid uint64) {
	limitTime := ja - 86400 * 7
	Model().Where("show_id = ? and job_at <= ?", sid, limitTime).Delete(nil)

	d := &Table{
		ShowId: sid,
		PlatformId: pid,
		Content: content,
		Author: name,
		PublishAt: _time,
		Forward: forward,
		JobAt: ja,
	}
	Model().Create(&d)
}