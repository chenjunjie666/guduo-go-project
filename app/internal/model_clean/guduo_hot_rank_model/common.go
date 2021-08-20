package guduo_hot_rank_model

import (
	"gorm.io/gorm"
	"guduo/app/crawler/clean/task"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前骨朵热度TOP
func SaveCurRank(hot []*task.GuduoHotItem, type_ int64, cycle int8, da uint, sid uint64) {
	res := make([]Table, len(hot))

	for k, v := range hot {
		res[k] = Table{
			ShowId:      sid,
			ShowType:    type_,
			SubShowType: v.SubType,
			Num:         v.Hot,
			DayAt:       da,
		}
	}

	if len(res) > 0 {
		Model().Create(&res)
	}
}