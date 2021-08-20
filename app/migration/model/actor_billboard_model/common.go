package actor_billboard_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/actor_domi_rank_model"
	"guduo/app/internal/model_clean/actor_hot_rank_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/migration/internal"
	"guduo/app/migration/model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetLoliPopMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync() {
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").
			Where("day >= ?", internal.StartDay).
			Offset(offset).
			Limit(limit).
			Find(&res)

		sync(res)

		if len(res) < 100000 {
			break
		}

		offset += limit
	}
}

func sync (res []*Table){
	arr := make([]*actor_hot_rank_model.Table, 0 ,500)
	arr2 := make([]*actor_domi_rank_model.Table, 0 ,500)
	for _, v := range res {
		type1, type2 := getType(v.Category)
		if type1 == -10 && type2 == -10 {
			continue
		}

		if type1 < 0 {
			isNew := int8(1)
			if type1 == -1 {
				isNew = 0
			}

			row := &actor_hot_rank_model.Table{
				ActorId:   v.ActorId,
				ActorName: v.ActorName,
				IsNew:     isNew,
				Cycle:     type2,
				Num:       v.EffectionIndex,
				CustomNum: 0,
				Rank:      v.BillboardRank,
				Rise:      v.EffectionRankRise,
				DayAt:     model.StrToTime(v.Day),
			}

			arr = append(arr, row)
			if len(arr) >= 400 {
				actor_hot_rank_model.Model().Create(&arr)
				arr = arr[:0]
			}
		} else {
			row := &actor_domi_rank_model.Table{
				ActorId:   v.ActorId,
				ActorName: v.ActorName,
				PlayType:  type1,
				Cycle:     type2,
				Num:       v.EffectionIndex,
				CustomNum: 0,
				Rank:      v.BillboardRank,
				Rise:      v.EffectionRankRise,
				DayAt:     model.StrToTime(v.Day),
			}

			arr2 = append(arr2, row)
			if len(arr2) >= 400 {
				actor_domi_rank_model.Model().Create(&arr2)
				arr2 = arr2[:0]
			}
		}
	}


	if len(arr) > 0 {
		actor_hot_rank_model.Model().Create(&arr)
	}
	if len(arr2) > 0 {
		actor_domi_rank_model.Model().Create(&arr2)
	}
}

func getType(c string) (int8, int8) {
	switch c {
	case "DRAMA_R1":
		return show_actor_model.PlayTypeLead, model_clean.CycleDaily
	case "DRAMA_R1_HALF_YEAR":
		return show_actor_model.PlayTypeLead, model_clean.CycleYearly
	case "DRAMA_R1_MONTHLY":
		return show_actor_model.PlayTypeLead, model_clean.CycleMonthly
	case "DRAMA_R1_WEEKLY":
		return show_actor_model.PlayTypeLead, model_clean.CycleWeekly
	case "DRAMA_R2":
		return show_actor_model.PlayTypeStar, model_clean.CycleDaily
	case "DRAMA_R2_HALF_YEAR":
		return show_actor_model.PlayTypeStar, model_clean.CycleYearly
	case "DRAMA_R2_MONTHLY":
		return show_actor_model.PlayTypeStar, model_clean.CycleMonthly
	case "DRAMA_R2_WEEKLY":
		return show_actor_model.PlayTypeStar, model_clean.CycleWeekly
	case "DRAMA_R3":
		return show_actor_model.PlayTypeSupp, model_clean.CycleDaily
	case "DRAMA_R3_HALF_YEAR":
		return show_actor_model.PlayTypeSupp, model_clean.CycleYearly
	case "DRAMA_R3_MONTHLY":
		return show_actor_model.PlayTypeSupp, model_clean.CycleMonthly
	case "DRAMA_R3_WEEKLY":
		return show_actor_model.PlayTypeSupp, model_clean.CycleWeekly
	case "EFFECTION_DAILY":
		return -1, model_clean.CycleDaily
	case "EFFECTION_MONTHLY":
		return -1, model_clean.CycleMonthly
	case "EFFECTION_WEEKLY":
		return -1, model_clean.CycleWeekly
	case "NETWORK_DRAMA_R1":
		return -10, -10
	case "NETWORK_DRAMA_R1_HALF_YEAR":
		return -10, -10
	case "NETWORK_DRAMA_R1_MONTHLY":
		return -10, -10
	case "NETWORK_DRAMA_R1_WEEKLY":
		return -10, -10
	case "NEWSTAR_DAILY":
		return -2, model_clean.CycleDaily
	case "NEWSTAR_MONTHLY":
		return -2, model_clean.CycleMonthly
	case "NEWSTAR_WEEKLY":
		return -2, model_clean.CycleWeekly
	case "TV_DRAMA_R1":
		return -10, -10
	case "TV_DRAMA_R1_HALF_YEAR":
		return -10, -10
	case "TV_DRAMA_R1_MONTHLY":
		return -10, -10
	case "TV_DRAMA_R1_WEEKLY":
		return -10, -10
	case "U25_DAILY":
		return -2, model_clean.CycleDaily
	case "U25_MONTHLY":
		return -2, model_clean.CycleMonthly
	case "U25_WEEKLY":
		return -2, model_clean.CycleWeekly
	}
	return -10, -10
}
