package half_month

import (
	"guduo/app/crawler/clean/task"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/actor_hot_daily_model"
	"guduo/app/internal/model_clean/actor_hot_monthly_model"
	"guduo/app/internal/model_clean/actor_hot_rank_model"
	"guduo/app/internal/model_scrawler/actor_model"
)

func guduoActorHotHandle() {
	var res []actor_hot_daily_model.Table

	mdl := actor_hot_daily_model.Model()

	mdl.Select("actor_id", "IF(custom_num != 0, custom_num, num) as num", "is_new").
		Where("day_at >= ?", StartAt).
		Where("day_at <= ?", EndAt).
		Find(&res)

	tmp := make(map[uint64]float64)
	tmpDay := make(map[uint64]int)
	tmpIsNew := make(map[uint64]int8)

	for _, row := range res {
		if _, ok := tmp[row.ActorId]; ! ok {
			tmp[row.ActorId] = 0
			tmpDay[row.ActorId] = 0
			tmpIsNew[row.ActorId] = 0
		}

		tmp[row.ActorId] += row.Num
		tmpDay[row.ActorId]++
		tmpIsNew[row.ActorId] = row.IsNew
	}

	allHot := make([]*task.ActHotItem, 0, 500)
	newHot := make([]*task.ActHotItem, 0, 500)
	save := make([]actor_hot_monthly_model.Table, len(tmp))
	k := 0
	for aid, hotSum := range tmp {

		name := actor_model.GetActorName(aid)
		save[k] = actor_hot_monthly_model.Table{
			ActorId: aid,
			Num: hotSum / float64(tmpDay[aid]),
			DayAt: StartAt,
		}

		allHot = append(allHot, &task.ActHotItem{aid, name, save[k].Num})
		if tmpIsNew[aid] == 1 {
			newHot = append(newHot, &task.ActHotItem{aid, name, save[k].Num})
		}

		k++
	}

	hot := task.ActorTopK(allHot, 150)
	hot2 := task.ActorTopK(newHot, 100)
	actor_hot_rank_model.SaveCurRank(hot, 0, model_clean.CycleMonthly, JobAt)
	actor_hot_rank_model.SaveCurRank(hot2, 1, model_clean.CycleMonthly, JobAt)


	actor_hot_monthly_model.SaveMonthlyHot(save, StartAt)
}