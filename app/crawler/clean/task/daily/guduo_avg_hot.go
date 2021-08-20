package daily

import (
	"guduo/app/internal/model_clean/guduo_hot_avg_rank_model"
	"guduo/app/internal/model_clean/guduo_hot_rank_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/util"
)

// 计算剧的平均热度与排名

func AvgGuduoHotHandle() {
	avgCalc()
}

type avgHot struct {
	Num float64
	Day int64
}

func avgCalc() {
	sids := show_model.GetActiveShows()
	d := make(map[uint64]float64)
	var hots []guduo_hot_rank_model.Table
	guduo_hot_rank_model.Model().Select("show_id", "num").Where("show_id IN ?", sids).Find(&hots)
	tmp := make(map[uint64]*avgHot)

	for _, v := range hots {
		if _, ok := tmp[v.ShowId]; !ok {
			tmp[v.ShowId] = &avgHot{
				Num:         0,
				Day:       0,
			}
		}

		tmp[v.ShowId].Num += v.Num
		tmp[v.ShowId].Day++
	}

	for sid, row := range tmp{
		avg := row.Num / float64(row.Day)
		d[sid] = util.ToFixedFloat(avg, 2)
	}

	rank := make([]guduo_hot_avg_rank_model.Table, 0, 1000)
	for sid, hot := range d {
		rank = append(rank, guduo_hot_avg_rank_model.Table{
			ShowId:    sid,
			Num:       hot,
			Rank:      0,
		})
	}

	for i := 0; i < len(rank); i++ {
		for j := i + 1; j < len(rank); j++ {
			if rank[i].Num < rank[j].Num {
				rank[i], rank[j] = rank[j], rank[i]
			}
		}
	}

	for idx := range rank {
		rank[idx].Rank = int64(idx + 1)
		rank[idx].DayAt = JobAt
	}

	guduo_hot_avg_rank_model.Model().Where("id > ?", 0).Delete(nil)

	insert := make([]guduo_hot_avg_rank_model.Table, 0, 500)

	for _, v := range rank {
		insert = append(insert, v)
		if len(insert) >= 400 {
			guduo_hot_avg_rank_model.Model().Create(&insert)
			insert = insert[:0]
		}
	}

	if len(insert) > 0 {
		guduo_hot_avg_rank_model.Model().Create(&insert)
	}

}