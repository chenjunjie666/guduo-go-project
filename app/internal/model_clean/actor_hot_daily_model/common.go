package actor_hot_daily_model

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


func GetHot(aid []uint64, da uint) []*Table {
	var row []*Table

	Model().Select("IF(custom_num, custom_num, num) as num", "actor_id").
		Where("actor_id IN ?", aid).
		Where("day_at", da).
		Find(&row)

	return row
}

// 保存当日当前骨朵热度
func SaveCurHot(hot float64, isNew int8,  da uint, aid uint64) {
	row := &Table{
		ActorId: aid,
		IsNew: isNew,
		Num: hot,
		DayAt: da,
	}

	r := Model().Where("actor_id = ? and day_at = ?", aid, da).Limit(1).Find(&row)
	if r.RowsAffected > 0 {
		r.Updates(row)
	}else{
		Model().Create(&row)
	}
}