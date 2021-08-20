package guduo_hot_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/db"
	"guduo/pkg/util"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}


func GetGuduoHotByActor(aid []uint64, type_ int8, da uint) map[uint64][]float64 {
	allowSids := show_model.GetOnShowing()


	sids := show_actor_model.GetShowIdsByType(aid, type_, allowSids)

	sidMap := make(map[uint64]bool)
	for _, row := range sids {
		sidMap[row.ShowId] = true
	}
	sidArr := make([]uint64, 0, 1000)
	for sid, _ := range sidMap{
		sidArr = append(sidArr, sid)
	}

	var res []*Table
	Model().Select("IF(custom_num != 0, custom_num, num) as num").
		Where("`show_id` IN ?", sidArr).
		Where("`day_at` = ?", da).
		Order("num desc").
		Find(&res)

	ret := make(map[uint64][]float64)
	for _, row := range sids {
		if _, ok := ret[row.ActorId]; !ok {
			ret[row.ActorId] = make([]float64, 0, 5) // 只取前10部热播的剧
		}
		if len(ret[row.ActorId]) > 5 {
			continue
		}
		for _, hot := range res {
			if hot.ShowId == row.ShowId {
				h := util.ToFixedFloat(hot.Num, 2)
				ret[row.ActorId] = append(ret[row.ActorId], h)
				break
			}
		}
	}

	return ret
}


// 保存当日当前骨朵热度
//func SaveCurHot(hot float64, da uint, sid uint64) {
//	row := &Table{
//		ShowId: sid,
//		Num: hot,
//		DayAt: da,
//	}
//
//	r := Model().Where("show_id = ? and day_at = ?", sid, da).Limit(1).Find(&row)
//	if r.RowsAffected > 0 {
//		r.Updates(row)
//	}else{
//		Model().Create(&row)
//	}
//}