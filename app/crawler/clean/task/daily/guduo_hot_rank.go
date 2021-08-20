package daily

import (
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/guduo_hot_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_rank_model"
	"guduo/app/internal/model_scrawler/show_model"
)

func GuduoHotDailyRankHandle() {
	guduoHotDailyRank()
}


func guduoHotDailyRank() {
	var res []guduo_hot_daily_model.Table

	mdl := guduo_hot_daily_model.Model()

	mdl.Select("show_id", "IF(custom_num != 0, custom_num, num) as num").
		Where("day_at", Yesterday).
		Find(&res)

	//fmt.Println(res)
	//return

	rank := make([]*guduo_hot_rank_model.Table, 0, 1000)

	sids := make([]uint64, 0, 100)
	for _, row := range res {
		sids = append(sids, row.ShowId)
		rank = append(rank, &guduo_hot_rank_model.Table{
			ShowId: row.ShowId,
			Num:    row.Num,
			DayAt:  Yesterday,
			Rank:   0,
		})
	}

	var tmp []show_model.Table
	show_model.Model().Select("id", "platform", "show_type", "sub_show_type").
		Where("id IN ?", sids).Find(&tmp)

	type r struct {
		PlatformId  []int
		ShowType    int64
		SubShowType int64
	}

	info := make(map[uint64]*r)
	for _, row := range tmp {
		info[row.ID] = &r{
			PlatformId:  show_model.GetPlatform(row.Platform),
			ShowType:    row.ShowType,
			SubShowType: row.SubShowType,
		}
	}

	rankFin := make(map[int64]map[int64]map[int][]*guduo_hot_rank_model.Table)

	for _, row := range rank {
		i := info[row.ShowId]
		if i == nil {
			continue
		}
		if _, ok := rankFin[i.ShowType]; !ok {
			rankFin[i.ShowType] = make(map[int64]map[int][]*guduo_hot_rank_model.Table)
		}

		// 全类型
		if _, ok := rankFin[i.ShowType][-1]; !ok {
			rankFin[i.ShowType][-1] = make(map[int][]*guduo_hot_rank_model.Table)
		}
		// 全类型全平台
		if _, ok := rankFin[i.ShowType][-1][0]; !ok {
			rankFin[i.ShowType][-1][0] = make([]*guduo_hot_rank_model.Table, 0, 500)
		}

		rankFin[i.ShowType][-1][0] = append(rankFin[i.ShowType][-1][0], row)

		// 子类型
		if _, ok := rankFin[i.ShowType][i.SubShowType]; !ok {
			rankFin[i.ShowType][i.SubShowType] = make(map[int][]*guduo_hot_rank_model.Table)
		}
		// 子类型全平台
		if _, ok := rankFin[i.ShowType][i.SubShowType][0]; !ok {
			rankFin[i.ShowType][i.SubShowType][0] = make([]*guduo_hot_rank_model.Table, 0, 500)
		}
		rankFin[i.ShowType][i.SubShowType][0] = append(rankFin[i.ShowType][i.SubShowType][0], row)

		for _, plt := range i.PlatformId {
			if _, ok := rankFin[i.ShowType][i.SubShowType][plt]; !ok {
				rankFin[i.ShowType][i.SubShowType][plt] = make([]*guduo_hot_rank_model.Table, 0, 500)
			}

			rankFin[i.ShowType][i.SubShowType][plt] = append(rankFin[i.ShowType][i.SubShowType][plt], row)
			rankFin[i.ShowType][-1][plt] = append(rankFin[i.ShowType][-1][plt], row)
		}
	}

	guduo_hot_rank_model.Model().
		Where("day_at", Yesterday).
		Where("rank_type", model_clean.CycleDaily).
		Delete(nil)

	for st, row1 := range rankFin {
		for sst, row2 := range row1 {
			for plt, ranks := range row2 {
				sids2 := make([]uint64, 0, 100)
				for _, row := range ranks {
					sids2 = append(sids2, row.ShowId)
				}

				var tmp2 []*guduo_hot_rank_model.Table
				guduo_hot_rank_model.Model().Select("show_id", "rank").Where("show_id IN ?", sids2).
					Where("show_type", st).
					Where("sub_show_type", sst).
					Where("platform_id", plt).
					Where("rank_type", model_clean.CycleDaily).
					Where("day_at", beforeYesterday).
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

				save := make([]*guduo_hot_rank_model.Table, 0, 500)
				for idx := range ranks {
					lastr := lastRank[ranks[idx].ShowId]
					rise := int64(0)
					if lastr != 0 {
						rise = lastr - int64(idx+1)
					}
					ranks[idx].ID = 0
					ranks[idx].Rank = int64(idx + 1)
					ranks[idx].DayAt = Yesterday
					ranks[idx].Rise = rise
					ranks[idx].ShowType = st
					ranks[idx].SubShowType = sst
					ranks[idx].PlatformId = int64(plt)
					ranks[idx].RankType = model_clean.CycleDaily

					save = append(save, ranks[idx])

					if len(save) >= 400 {
						guduo_hot_rank_model.Model().Create(&save)
						save = save[:0]
					}
				}

				if len(save) > 0 {
					guduo_hot_rank_model.Model().Create(&save)
				}
			}
		}
	}
}










	//
	//
	//
	//
	//
	//tmp := make(map[uint64]float64)
	//tmpDay := make(map[uint64]int)
	//
	//for _, row := range res {
	//	if _, ok := tmp[row.ShowId]; !ok {
	//		tmp[row.ShowId] = 0
	//		tmpDay[row.ShowId] = 0
	//	}
	//
	//	tmp[row.ShowId] += row.Num
	//	tmpDay[row.ShowId]++
	//}
	//
	//save := make([]guduo_hot_monthly_model.Table, len(tmp))
	//k := 0
	//for sid, hotSum := range tmp {
	//	save[k] = guduo_hot_monthly_model.Table{
	//		ShowId: sid,
	//		Num:    util.ToFixedFloat(hotSum/float64(tmpDay[sid]), 2)
	//		DayAt:  StartAt,
	//	}
	//	k++
	//}
	//
	//for i := 0; i < len(save); i++ {
	//	for j := i; j < len(save); j++ {
	//		if save[j].Num > save[i].Num {
	//			save[i], save[j] = save[j], save[i]
	//		}
	//	}
	//}
	//
	//hotRank := make([]*guduo_hot_rank_model.Table, 0, 400)
	//for _, row := range save {
	//	hotRank = append(hotRank, &guduo_hot_rank_model.Table{
	//		ShowId:      0,
	//		ShowType:    0,
	//		SubShowType: 0,
	//		PlatformId:  0,
	//		Num:         0,
	//		CustomNum:   0,
	//		DayAt:       0,
	//	})
	//}
	//
	//allHot := make(map[int64]*task.GuduoHotRank)
	//for _, v := range save {
	//	row := show_model.GetShowsType(v.ShowId)
	//	if _, ok := allHot[row.ShowType]; !ok {
	//		allHot[row.ShowType] = &task.GuduoHotRank{
	//			Sid:  row.ID,
	//			Type: row.ShowType,
	//			Num:  make([]*task.GuduoHotItem, 0, 300),
	//		}
	//	}
	//	allHot[row.ShowType].Num = append(allHot[row.ShowType].Num, &task.GuduoHotItem{row.SubShowType, v.Num})
	//}
	//
	////for
	//
	//for type_, row := range allHot {
	//	top150 := task.GuduoTopK(row.Num, 150)
	//	guduo_hot_rank_model.SaveCurRank(top150, type_, model_clean.CycleMonthly, JobAt, row.Sid)
	//}
	//
	//guduo_hot_monthly_model.SaveMonthlyHot(save, StartAt)