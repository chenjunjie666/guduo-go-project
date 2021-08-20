package daily

import (
	"guduo/app/internal/model_clean/play_count_total_daily_model"
	"guduo/app/internal/model_clean/play_count_trend_daily_model"
	"guduo/app/internal/model_scrawler/show_model"
)

func PlayCountTotalHandle() {
	calcTotalPlayCount()
}

func calcTotalPlayCount(){
	// 总榜不需要区分类型，因为没有排名升降，所以小程序段，查找对应的类型show_id就行
	var d []play_count_trend_daily_model.Table
	play_count_trend_daily_model.Model().Select("show_id", "sum(num) as num").
		Where("day_at < ?", JobAt).
		Group("show_id").
		Find(&d)


	rank := make([]*play_count_total_daily_model.Table, 0, 1000)

	sids := make([]uint64, 0, 100)
	for _, row := range d {
		sids = append(sids, row.ShowId)
		rank = append(rank, &play_count_total_daily_model.Table{
			ShowId: row.ShowId,
			Num:    row.Num,
			DayAt:  Yesterday,
			Rank: 0,
		})
	}

	//fmt.Println(rank)

	var tmp []show_model.Table
	show_model.Model().Select("id", "platform", "show_type", "sub_show_type").Where("id IN ?", sids).Find(&tmp)

	type r struct {
		ShowType int64
		SubShowType int64
	}

	info := make(map[uint64]*r)
	for _, row := range tmp {
		info[row.ID] = &r{
			ShowType:    row.ShowType,
			SubShowType: row.SubShowType,
		}
	}

	rankFin := make(map[int64]map[int64][]*play_count_total_daily_model.Table)

	for _, row := range rank {
		i := info[row.ShowId]
		if i == nil {
			continue
		}
		if _, ok := rankFin[i.ShowType]; !ok {
			rankFin[i.ShowType] = make(map[int64][]*play_count_total_daily_model.Table)
		}

		// 全类型
		if _, ok := rankFin[i.ShowType][-1]; !ok {
			rankFin[i.ShowType][-1] = make([]*play_count_total_daily_model.Table, 0, 500)
		}

		rankFin[i.ShowType][-1] = append(rankFin[i.ShowType][-1], row)

		// 子类型
		if _, ok := rankFin[i.ShowType][i.SubShowType]; !ok {
			rankFin[i.ShowType][i.SubShowType] = make([]*play_count_total_daily_model.Table, 0, 500)
		}

		rankFin[i.ShowType][i.SubShowType] = append(rankFin[i.ShowType][i.SubShowType], row)
	}

	play_count_total_daily_model.Model().Where("day_at", Yesterday).Delete(nil)

	for st, row1 := range rankFin {
		for sst, ranks := range row1 {
			sids2 := make([]uint64, 0, 100)
			for _, row := range ranks {
				sids2 = append(sids2, row.ShowId)
			}

			var tmp2 []*play_count_total_daily_model.Table
			play_count_total_daily_model.Model().Select("show_id", "rank").Where("show_id IN ?", sids2).
				Where("show_type", st).
				Where("sub_show_type", sst).
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

			save := make([]*play_count_total_daily_model.Table, 0, 500)
			for idx := range ranks {
				lastr := lastRank[ranks[idx].ShowId]
				rise := int64(0)
				if lastr != 0 {
					rise = lastr - int64(idx + 1)
				}
				ranks[idx].ID = 0
				ranks[idx].Rank = int64(idx + 1)
				ranks[idx].DayAt = Yesterday
				ranks[idx].Rise = rise
				ranks[idx].ShowType = st
				ranks[idx].SubShowType = sst

				save = append(save, ranks[idx])

				if len(save) >= 400 {
					play_count_total_daily_model.Model().Create(&save)
					save = save[:0]
				}
			}

			if len(save) > 0 {
				play_count_total_daily_model.Model().Create(&save)
			}
		}
	}
















	//for i := 0; i < len(rank); i++ {
	//	for j := i + 1; j < len(rank); j++ {
	//		if rank[i].Num < rank[j].Num {
	//			rank[i], rank[j] = rank[j], rank[i]
	//		}
	//	}
	//}
	//
	//for idx := range rank {
	//	rank[idx].DayAt = JobAt
	//}
	//
	//play_count_total_daily_model.Model().Where("day_at", JobAt).Delete(nil)
	//insert := make([]*play_count_total_daily_model.Table, 0, 500)
	//for _, v := range rank {
	//	insert = append(insert, v)
	//	if len(insert) >= 400 {
	//		play_count_total_daily_model.Model().Create(&insert)
	//		insert = insert[:0]
	//	}
	//}
	//
	//if len(insert) > 0 {
	//	play_count_total_daily_model.Model().Create(&insert)
	//}
}