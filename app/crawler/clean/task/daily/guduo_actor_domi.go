package daily

import (
	"guduo/app/crawler/clean/task"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/actor_domi_daily_model"
	"guduo/app/internal/model_clean/actor_domi_rank_model"
	"guduo/app/internal/model_clean/actor_hot_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_daily_model"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/pkg/time"
	"guduo/pkg/util"
)

func GuduoActorDomiHandle() {
	guduoActorDomi()
}

// 骨朵艺人霸屏榜
func guduoActorDomi() {

	hot1 := make([]*task.ActHotItem, 0, 100)
	hot2 := make([]*task.ActHotItem, 0, 100)
	hot3 := make([]*task.ActHotItem, 0, 100)
	data := guduoActorDomiGetInfo()
	dailyHot := make([]*actor_domi_daily_model.Table, 0, 10000)
	for _, row := range data {
		aid := row.Aid

		avgDomiHot := .0
		if len(row.GuduoHot) > 0 {
			for _, v := range row.GuduoHot {
				avgDomiHot += v
			}
			avgDomiHot = avgDomiHot / float64(len(row.GuduoHot))
		}

		domi := row.ActorHot*0.4 + avgDomiHot*0.6

		if row.Type == show_actor_model.PlayTypeLead{
			hot1 = append(hot1, &task.ActHotItem{row.Aid, row.Name, domi})
		}else if row.Type == show_actor_model.PlayTypeStar {
			hot2 = append(hot2, &task.ActHotItem{row.Aid, row.Name, domi})
		}else {
			hot3 = append(hot3, &task.ActHotItem{row.Aid, row.Name, domi})
		}

		dailyHot = append(dailyHot, &actor_domi_daily_model.Table{
			ActorId:   aid,
			Num:       domi,
			Type:      row.Type,
			DayAt:     JobAt - 86400,
		})

	}
	actor_domi_daily_model.SaveCurHot(dailyHot, JobAt - 86400)

	actor_domi_rank_model.SaveCurRank(hot1, show_actor_model.PlayTypeLead, model_clean.CycleDaily, JobAt - 86400)
	actor_domi_rank_model.SaveCurRank(hot2, show_actor_model.PlayTypeStar, model_clean.CycleDaily, JobAt - 86400)
	actor_domi_rank_model.SaveCurRank(hot3, show_actor_model.PlayTypeSupp, model_clean.CycleDaily, JobAt - 86400)
}

func guduoActorDomiGetInfo() []*task.GuduoActorDomiIndicator {
	day := time.Today() - 86400
	aids := actor_model.GetActor()

	data := make([]*task.GuduoActorDomiIndicator, 0, 100)
	idx := 0

	aidArr := make([]uint64, len(aids))

	for k, v := range aids {
		aidArr[k] = v.Id
	}


	actorsHot := actor_hot_daily_model.GetHot(aidArr, day)
	guduoHotLead := guduo_hot_daily_model.GetGuduoHotByActor(aidArr, show_actor_model.PlayTypeLead, day)
	guduoHotStar := guduo_hot_daily_model.GetGuduoHotByActor(aidArr, show_actor_model.PlayTypeStar, day)
	guduoHotSupp := guduo_hot_daily_model.GetGuduoHotByActor(aidArr, show_actor_model.PlayTypeSupp, day)

	for _, aid := range aids {
		row := &task.GuduoActorDomiIndicator{
			Aid:      aid.Id,
			Name:     aid.Name,
			Type:     show_actor_model.PlayTypeLead,
			ActorHot: .0,
			GuduoHot: make([]float64, 0, 10),
		}

		row2 := &task.GuduoActorDomiIndicator{
			Aid:      aid.Id,
			Name:     aid.Name,
			Type:     show_actor_model.PlayTypeStar,
			ActorHot: .0,
			GuduoHot: make([]float64, 0, 10),
		}

		row3 := &task.GuduoActorDomiIndicator{
			Aid:      aid.Id,
			Name:     aid.Name,
			Type:     show_actor_model.PlayTypeSupp,
			ActorHot: .0,
			GuduoHot: make([]float64, 0, 10),
		}
		for _, item := range actorsHot {
			if aid.Id == item.ActorId {
				row.ActorHot = util.ToFixedFloat(item.Num, 2)
				row2.ActorHot = util.ToFixedFloat(item.Num, 2)
				row3.ActorHot = util.ToFixedFloat(item.Num, 2)
				break
			}
		}

		for actorId, item := range guduoHotLead {
			if aid.Id == actorId {
				row.GuduoHot = item
				break
			}
		}
		for actorId, item := range guduoHotStar {
			if aid.Id == actorId {
				row2.GuduoHot = item
				break
			}
		}
		for actorId, item := range guduoHotSupp {
			if aid.Id == actorId {
				row3.GuduoHot = item
				break
			}
		}
		idx++
		data = append(data, row)  // 领衔主演
		data = append(data, row2) // 主演
		data = append(data, row3) // 配角
	}

	return data
}
