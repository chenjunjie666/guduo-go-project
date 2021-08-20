package daily

import (
	"guduo/app/internal/model_clean/play_count_trend_daily_model"
	"guduo/app/internal/model_clean/play_count_yearly_model"
	"guduo/app/internal/model_scrawler/show_model"
	"time"
)

func YearTotalPlayCountHandle(){
	calcYearPlayCount()
}

func calcYearPlayCount() {
	var d []play_count_trend_daily_model.Table
	play_count_trend_daily_model.Model().Select("show_id", "sum(num) as num").
		Where("day_at between ? and ?", yearStart, JobAt).
		Group("show_id").
		Find(&d)


	rank := make([]*play_count_yearly_model.Table, 0, 1000)

	sids := make([]uint64, 0, 100)
	for _, row := range d {
		sids = append(sids, row.ShowId)
		rank = append(rank, &play_count_yearly_model.Table{
			ShowId: row.ShowId,
			Num:    row.Num,
			DayAt:  yearStart,
			Rank:   0,
		})
	}

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

	rankFin := make(map[int64]map[int64][]*play_count_yearly_model.Table)

	for _, row := range rank {
		i := info[row.ShowId]
		if i == nil {
			continue
		}
		if _, ok := rankFin[i.ShowType]; !ok {
			rankFin[i.ShowType] = make(map[int64][]*play_count_yearly_model.Table)
		}

		// 全类型
		if _, ok := rankFin[i.ShowType][-1]; !ok {
			rankFin[i.ShowType][-1] = make([]*play_count_yearly_model.Table, 0, 500)
		}

		rankFin[i.ShowType][-1] = append(rankFin[i.ShowType][-1], row)

		// 子类型
		if _, ok := rankFin[i.ShowType][i.SubShowType]; !ok {
			rankFin[i.ShowType][i.SubShowType] = make([]*play_count_yearly_model.Table, 0, 500)
		}

		rankFin[i.ShowType][i.SubShowType] = append(rankFin[i.ShowType][i.SubShowType], row)
	}

	play_count_yearly_model.Model().Where("day_at", yearStart).Delete(nil)
	lastYear := time.Unix(int64(yearStart), 0).AddDate(-1, 0, 0).Unix()

	for st, row1 := range rankFin {
		for sst, ranks := range row1 {
			sids2 := make([]uint64, 0, 100)
			for _, row := range ranks {
				sids2 = append(sids2, row.ShowId)
			}

			var tmp2 []*play_count_yearly_model.Table
			play_count_yearly_model.Model().Select("show_id", "rank").Where("show_id IN ?", sids2).
				Where("show_type", st).
				Where("sub_show_type", sst).
				Where("day_at", lastYear).
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

			save := make([]*play_count_yearly_model.Table, 0, 500)
			for idx := range ranks {
				lastr := lastRank[ranks[idx].ShowId]
				rise := int64(0)
				if lastr != 0 {
					rise = lastr - int64(idx + 1)
				}
				ranks[idx].ID = 0
				ranks[idx].Rank = int64(idx + 1)
				ranks[idx].DayAt = yearStart
				ranks[idx].Rise = rise
				ranks[idx].ShowType = st
				ranks[idx].SubShowType = sst

				save = append(save, ranks[idx])

				if len(save) >= 400 {
					play_count_yearly_model.Model().Create(&save)
					save = save[:0]
				}
			}

			if len(save) > 0 {
				play_count_yearly_model.Model().Create(&save)
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
	//	rank[idx].Rank = int64(idx + 1)
	//	rank[idx].DayAt = StartAt
	//}
	//
	//play_count_yearly_model.Model().Where("day_at", StartAt).Delete(nil)
	//play_count_yearly_model.Model().Create(&rank)
}


//var tmp []show_model.Table
//show_model.Model().Select("id", "platform", "show_type", "sub_show_type").Where("id IN ?", sids).Find(&tmp)
//
//type r struct {
//	PlatformId []int
//	ShowType int64
//	SubShowType int64
//}
//
//info := make(map[uint64]*r)
//for _, row := range tmp {
//p := show_model.GetPlatform(row.Platform)
//info[row.ID] = &r{
//PlatformId:  p,
//ShowType:    row.ShowType,
//SubShowType: row.SubShowType,
//}
//}


//rankFin := make(map[int64]map[int64]map[int][]*play_count_yearly_model.Table)
//
//for _, row := range rank {
//i := info[row.ShowId]
//
//if _, ok := rankFin[i.ShowType]; !ok {
//rankFin[i.ShowType] = make(map[int64]map[uint64][]*play_count_yearly_model.Table)
//}
//
//// 全类型
//if _, ok := rankFin[i.ShowType][-1]; !ok {
//rankFin[i.ShowType][-1] = make(map[int][]*play_count_yearly_model.Table)
//}
//
//// 全平台
//if _, ok := rankFin[i.ShowType][-1][0]; !ok {
//rankFin[i.ShowType][-1][0] = make([]*play_count_yearly_model.Table, 0, 500)
//}
//rankFin[i.ShowType][-1][0] = append(rankFin[i.ShowType][-1][0], row)
//
//// 子类型
//if _, ok := rankFin[i.ShowType][i.SubShowType]; !ok {
//rankFin[i.ShowType][i.SubShowType] = make(map[int][]*play_count_yearly_model.Table)
//}
//
//// 子类型全平台
//if _, ok := rankFin[i.ShowType][-1][0]; !ok {
//rankFin[i.ShowType][i.SubShowType][0] = make([]*play_count_yearly_model.Table, 0, 500)
//}
//rankFin[i.ShowType][i.SubShowType][0] = append(rankFin[i.ShowType][i.SubShowType][0], row)
//
//for _, plt := range i.PlatformId {
//if _, ok := rankFin[i.ShowType][-1][plt]; !ok {
//rankFin[i.ShowType][-1][plt] = make([]*play_count_yearly_model.Table, 0, 500)
//}
//rankFin[i.ShowType][-1][plt] = append(rankFin[i.ShowType][-1][plt], row)
//
//if _, ok := rankFin[i.ShowType][i.SubShowType][plt]; !ok {
//rankFin[i.ShowType][i.SubShowType][plt] = make([]*play_count_yearly_model.Table, 0, 500)
//}
//rankFin[i.ShowType][i.SubShowType][plt] = append(rankFin[i.ShowType][i.SubShowType][plt], row)
//}
//}
//
//
//for st, row1 := range rankFin {
//for sst, row2 := range row1 {
//for plt, row := range row2 {
//play_count_yearly_model.Model().Select()
//
//
//
//}
//}
//}