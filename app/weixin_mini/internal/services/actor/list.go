package actor

import (
	"fmt"
	"guduo/app/internal/model_clean/actor_domi_rank_model"
	"guduo/app/internal/model_clean/actor_hot_rank_model"
	"guduo/app/internal/model_clean/guduo_hot_rank_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_model"
)

const (
	ListTypeDay = iota
	ListTypeWeek
	ListTypeMonth
)

const (
	rankTypeHot = iota  // 热力榜
	rankTypeDomi // 霸屏榜
	rankTypeNewHot //新秀榜
)

func List(dayAt, listType, rankType int, playType string) []map[string]interface{} {
	switch rankType {
	case rankTypeHot:
		return hotRank(dayAt, listType)
	case rankTypeDomi:
		return domiRank(dayAt, listType, playType)
	case rankTypeNewHot:
		return newHotRank(dayAt, listType)
	}
	return nil
}

type hot struct {
	ActorId uint64
	Name string
	Num float64
}

func hotRank(dayAt int, lt int) []map[string]interface{} {
	mdl := actor_hot_rank_model.Model()
	var res []hot

	mdl.Select("actor_id", "actor_name as name", "IF(custom_num, custom_num, num) as num").
		Where("day_at", dayAt).
		Where("cycle", lt + 1).
		Where("is_new", 0).
		Order("num desc").
		Limit(50).
		Find(&res)

	ret := make([]map[string]interface{}, len(res))
	for k, row := range res {
		ret[k] = map[string]interface{}{
			"actor_id": row.ActorId,
			"hot": row.Num,
			"name": row.Name,
			"shows": nil,
		}
	}

	return ret
}

func newHotRank(dayAt int, lt int) []map[string]interface{} {
	mdl := actor_hot_rank_model.Model()
	var res []hot

	mdl.Select("actor_id", "actor_name as name", "IF(custom_num, custom_num, num) as num").
		Where("day_at", dayAt).
		Where("cycle", lt + 1).
		Where("is_new", 1).
		Order("num desc").
		Limit(50).
		Find(&res)

	ret := make([]map[string]interface{}, len(res))
	for k, row := range res {
		ret[k] = map[string]interface{}{
			"actor_id": row.ActorId,
			"hot": row.Num,
			"name": row.Name,
			"shows": nil,
		}
	}

	return ret
}


func domiRank(dayAt, lt int, playType string) []map[string]interface{} {
	mdl := actor_domi_rank_model.Model()
	var res []hot

	mdl.Select("actor_id", "actor_name as name", "IF(custom_num, custom_num, num) as num").
		Where("day_at", dayAt).
		Where("cycle", lt + 1).
		Where("play_type", playType).
		Order("num desc").
		Limit(50).
		Find(&res)

	aids := make([]uint64, len(res))
	ret := make([]map[string]interface{}, len(res))
	for k, row := range res {
		aids[k] = row.ActorId
		ret[k] = map[string]interface{}{
			"actor_id": row.ActorId,
			"hot": row.Num,
			"name": row.Name,
			"shows": "",
		}
	}

	//var actors []actor_model.Table
	//mdl = actor_model.Model()
	//mdl.Select("id", "name").
	//	Where("id IN ?", aids).
	//	Find(&actors)

	shows := guduoHotToActor(dayAt, lt)
	for k, row := range ret {
		aid, _ := row["actor_id"].(uint64)

		for aid_, v := range shows {
			if aid == aid_ {
				ret[k]["shows"] = v
			}
		}
	}

	return ret
}

// 找出演员对应的热播剧
func guduoHotToActor(dayAt int, listType int) map[uint64][]string {
	mdl := guduo_hot_rank_model.Model()

	var sidRes []guduo_hot_rank_model.Table
	mdl.Select("show_id").
		Where("day_at", dayAt).
		Where("rank_type", listType + 1).
		Where("platform_id", 0).
		Where("show_type", show_model.ShowTypeSeries).
		Where("sub_show_type", -1).
		Order("IF(custom_num != 0, custom_num, num) desc").
		Limit(100).
		Find(&sidRes)

	sids := make([]uint64, len(sidRes))
	for k, row := range sidRes {
		sids[k] = row.ShowId
	}

	var res []show_actor_model.Table
	mdl = show_actor_model.Model()
	mdl.Select("show_id", "actor_id").
		Where("show_id IN ?", sids).
		Find(&res)

	var res2 []show_model.Table
	mdl = show_model.Model()
	mdl.Select("id", "name").Where("id IN ?", sids).
		Where("status", show_model.ShowStatStandard).
		Where("is_show", show_model.ShowOn).
		Find(&res2)

	ret := make(map[uint64][]string)
	tmp := make(map[uint64][]uint64)

	for _, v := range res {
		if _, ok := tmp[v.ActorId]; !ok {
			tmp[v.ActorId] = make([]uint64, 0, 2)
			ret[v.ActorId] = make([]string, 0, 2)
		}
		// 一个演员最多对应两个热门剧
		if len(tmp[v.ActorId]) < 2{
			if len(tmp[v.ActorId]) == 1 && tmp[v.ActorId][0] == v.ShowId {
				continue
			}
			tmp[v.ActorId] = append(tmp[v.ActorId], v.ShowId)
		}
	}

	fmt.Println(tmp)

	for aid, sids2 := range tmp {
		for _, row := range res2 {
			for _, sid := range sids2 {
				if sid == row.ID {
					ret[aid] = append(ret[aid], row.Name)
				}
			}
		}
	}

	return ret
}