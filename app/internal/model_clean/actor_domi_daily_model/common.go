package actor_domi_daily_model

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

// 保存当日当前艺人霸屏
func SaveCurHot(hot []*Table, da uint) {
	res := make([]*Table, 0, 400)

	Model().Where("day_at", da).Delete(nil)

	for _, v := range hot {
		res = append(res, v)

		if len(res) >= 400 {
			Model().Create(&res)
			res = res[:0]
		}
	}

	if len(res) > 0 {
		Model().Create(&res)
	}
}