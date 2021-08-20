package danmaku_word_cloud_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前弹幕词云
func SaveWordCloud(wc map[string]int, ja uint, sid uint64) {
	da := helper.JobAt2DayAt(ja)

	wcs := make([]*Table, 0, 50)
	for word, weight := range wc {
		wcs = append(wcs, &Table{
			ShowId: sid,
			Word:   word,
			Weight: int64(weight),
			DayAt:  da,
		})
	}

	Model().Where("show_id = ? and day_at = ?", sid, da).Delete(nil)
	Model().Create(&wcs)
}
