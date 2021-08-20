package half_month

import (
	"guduo/app/crawler/clean/task"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/actor_domi_daily_model"
	"guduo/app/internal/model_clean/actor_domi_monthly_model"
	"guduo/app/internal/model_clean/actor_domi_rank_model"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
)

func guduoActorDomiHandle() {
	var res []actor_domi_daily_model.Table

	mdl := actor_domi_daily_model.Model()

	mdl.Select("actor_id", "IF(custom_num != 0, custom_num, num) as num", "type").
		Where("day_at >= ?", StartAt).
		Where("day_at <= ?", EndAt).
		Find(&res)

	tmp := make(map[uint64]map[int8]float64)
	tmpDay := make(map[uint64]map[int8]int)

	for _, row := range res {
		if _, ok := tmp[row.ActorId]; ! ok {
			tmp[row.ActorId] = make(map[int8]float64)
			tmpDay[row.ActorId] = make(map[int8]int)
		}

		if _, ok := tmp[row.ActorId][row.Type]; ! ok {
			tmp[row.ActorId][row.Type] = 0
			tmpDay[row.ActorId][row.Type] = 0
		}
		tmp[row.ActorId][row.Type] += row.Num
		tmpDay[row.ActorId][row.Type]++
	}

	hot1 := make([]*task.ActHotItem, 0, 100)
	hot2 := make([]*task.ActHotItem, 0, 100)
	hot3 := make([]*task.ActHotItem, 0, 100)
	save := make([]actor_domi_monthly_model.Table, len(tmp), len(tmp) * 2)
	k := 0
	for aid, row := range tmp {
		name := actor_model.GetActorName(aid)
		for type_, hot := range row {
			save[k] = actor_domi_monthly_model.Table{
				ActorId: aid,
				Num: hot / float64(tmpDay[aid][type_]),
				DayAt: StartAt,
				Type: type_,
			}
			k++

			if type_ == show_actor_model.PlayTypeLead{
				hot1 = append(hot1, &task.ActHotItem{aid, name, save[k].Num})
			}else if type_ == show_actor_model.PlayTypeStar {
				hot2 = append(hot2, &task.ActHotItem{aid, name, save[k].Num})
			}else {
				hot3 = append(hot3, &task.ActHotItem{aid, name, save[k].Num})
			}
		}
	}

	actor_domi_rank_model.SaveCurRank(hot1, show_actor_model.PlayTypeLead, model_clean.CycleMonthly, JobAt)
	actor_domi_rank_model.SaveCurRank(hot2, show_actor_model.PlayTypeStar, model_clean.CycleMonthly, JobAt)
	actor_domi_rank_model.SaveCurRank(hot3, show_actor_model.PlayTypeSupp, model_clean.CycleMonthly, JobAt)

	actor_domi_monthly_model.SaveMonthlyHot(save, StartAt)
}