package actor_hot_rank_model

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

// 保存当日当前艺人热度
func SaveCurRank(hot []*task.ActHotItem, isNew int8, cycle int8, da uint) {
	res := make([]*Table, 0, 100)

	for i := 0; i <= len(hot); i++ {
		for j := i; j < len(hot); j++ {
			if hot[j].Hot > hot[i].Hot {
				hot[i], hot[j] = hot[j], hot[i]
			}
		}
	}
	for k, v := range hot {
		res = append(res, &Table{
			ActorId:   v.Aid,
			ActorName: v.ActorName,
			IsNew:     isNew,
			Cycle:     cycle,
			Num:       v.Hot,
			CustomNum: 0,
			Rank:      int64(k + 1),
			Rise:      0, // 排名升降不需要了
			DayAt:     da,
		})

		if len(res) >= 400 {
			Model().Create(&res)
			res = res[:0]
		}
	}

	if len(res) > 0 {
		Model().Create(&res)
	}

}