package half_week

import (
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/play_count_rank_model"
	"guduo/app/internal/model_clean/play_count_trend_daily_model"
	"guduo/app/internal/model_scrawler/show_model"
)

func moviePlayCountHandle() {

	sids := show_model.GetActiveShowsByType(show_model.ShowTypeMovie, show_model.ShowSubTypeMovieCinema)

	var res []play_count_trend_daily_model.Table

	mdl := play_count_trend_daily_model.Model()

	mdl.Select("show_id", "platform_id", "IF(custom_num != 0, custom_num, num) as num").
		Where("day_at >= ?", StartAt).
		Where("day_at <= ?", EndAt).
		Where("show_id IN ?", sids).
		Find(&res)

	//fmt.Println(res)
	//return

	rankFin := make(map[uint64][]*play_count_rank_model.Table)

	for _, v := range res {
		if _, ok := rankFin[v.PlatformId]; !ok {
			rankFin[v.PlatformId] = make([]*play_count_rank_model.Table, 0, 500)
		}

		find := false
		for k, row := range rankFin[v.PlatformId]{
			if row.ShowId == v.ShowId {
				find = true
				rankFin[v.PlatformId][k].Num += v.Num
				break
			}
		}
		if !find {
			rankFin[v.PlatformId] = append(rankFin[v.PlatformId], &play_count_rank_model.Table{
				ShowId:      v.ShowId,
				Num:         v.Num,
			})
		}
	}

	play_count_rank_model.Model().
		Where("day_at", StartAt).
		Where("rank_type", model_clean.CycleWeekly).
		Delete(nil)

	for plt, ranks := range rankFin {
		sids2 := make([]uint64, 0, 100)
		for _, row := range ranks {
			sids2 = append(sids2, row.ShowId)
		}

		var tmp2 []*play_count_rank_model.Table
		play_count_rank_model.Model().Select("show_id", "rank").Where("show_id IN ?", sids2).
			Where("show_type", 2).
			Where("sub_show_type", 21).
			Where("platform_id", plt).
			Where("rank_type", model_clean.CycleWeekly).
			Where("day_at", lastWeek).
			Find(&tmp2)

		lastRank := make(map[uint64]int64)
		for _, v := range tmp2 {
			lastRank[v.ShowId] = v.Rank
		}

		for i := 0; i < len(ranks); i++ {
			for j := i + 1; j < len(ranks); j++ {
				if ranks[i].Num < ranks[j].Num {
					ranks[i], ranks[j] = ranks[j], ranks[i]
				}
			}
		}

		save := make([]*play_count_rank_model.Table, 0, 500)
		for idx := range ranks {
			lastr := lastRank[ranks[idx].ShowId]
			rise := int64(0)
			if lastr != 0 {
				rise = lastr - int64(idx+1)
			}
			ranks[idx].ID = 0
			ranks[idx].Rank = int64(idx + 1)
			ranks[idx].DayAt = StartAt
			ranks[idx].Rise = rise
			ranks[idx].ShowType = 2
			ranks[idx].SubShowType = 21
			ranks[idx].PlatformId = plt
			ranks[idx].RankType = model_clean.CycleWeekly

			save = append(save, ranks[idx])

			if len(save) >= 400 {
				play_count_rank_model.Model().Create(&save)
				save = save[:0]
			}
		}

		if len(save) > 0 {
			play_count_rank_model.Model().Create(&save)
		}
	}
}
