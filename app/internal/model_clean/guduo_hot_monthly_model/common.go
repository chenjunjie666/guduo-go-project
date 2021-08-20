package guduo_hot_monthly_model

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

// 保存当周当前骨朵热度
func SaveMonthlyHot(data []Table, da uint) {
	Model().Where("day_at", da).Delete(nil)
	Model().Create(&data)
}