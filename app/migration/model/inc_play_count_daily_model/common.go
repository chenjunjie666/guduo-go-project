package inc_play_count_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/play_count_rank_model"
	play_count_trend_daily_model "guduo/app/internal/model_clean/play_count_trend_daily_model"
	"guduo/app/internal/model_clean/play_count_yearly_model"
	"guduo/app/internal/model_scrawler/show_model"
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
	daily_trend()
	play_rank()
	play_count_year()
}

func daily_trend() {
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").Where("day >= ?", internal.StartDay).
			Where("type", "DAILY").
			Offset(offset).
			Limit(limit).
			Find(&res)

		sync_daily_trend(res)

		if len(res) < 100000 {
			break
		}

		offset += limit
	}
}

func sync_daily_trend(res []*Table) {

	arr := make([]*play_count_trend_daily_model.Table, 0, 400)
	for _, v := range res {
		st, sst := getShowType(v.Category)
		if st == -10 {
			continue
		}
		if sst != -1 {
			continue
		}

		if v.PlatformId == 0 {
			continue
		}
		row := &play_count_trend_daily_model.Table{
			ShowId:     v.ShowId,
			PlatformId: v.PlatformId,
			Num:        v.PlayCount,
			DayAt:      model.StrToTime(v.Day),
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			play_count_trend_daily_model.Model().Create(&arr)
			arr = make([]*play_count_trend_daily_model.Table, 0, 400)
		}
	}
	if len(arr) > 0 {
		play_count_trend_daily_model.Model().Create(&arr)
	}
}

func play_rank() {
	offset := 0
	limit := 100000
	for {
		// 只同步院线电影
		var res []*Table
		Model().Select("*").
			Where("day >= ?", internal.StartDay).
			Where("type IN ?", []string{"DAILY", "WEEKLY", "MONTHLY"}).
			Where("category", "MOVIE").
			Offset(offset).
			Limit(limit).
			Find(&res)

		sync_play_rank(res)

		if len(res) < 100000 {
			break
		}

		offset += limit
	}
}

func sync_play_rank(res []*Table) {

	arr := make([]*play_count_rank_model.Table, 0, 400)
	for _, v := range res {
		st, sst := int64(2), int64(21)
		rankType := getType(v.Category)
		row := &play_count_rank_model.Table{
			ShowId:      v.ShowId,
			ShowType:    st,
			SubShowType: sst,
			PlatformId:  v.PlatformId,
			RankType:    rankType,
			Rank:        v.PlayCountRank,
			Rise:        v.PlayCountRise,
			Num:         v.PlayCount,
			DayAt:       model.StrToTime(v.Day),
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			play_count_rank_model.Model().Create(&arr)
			arr = make([]*play_count_rank_model.Table, 0, 400)
		}
	}
	if len(arr) > 0 {
		play_count_rank_model.Model().Create(&arr)
	}
}

func play_count_year() {
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").
			Where("day >= ?", internal.StartDay).
			Where("type", "YEARLY").
			Offset(offset).
			Limit(limit).
			Find(&res)

		sync_play_count_year(res)

		if len(res) < 100000 {
			break
		}

		offset += limit
	}

}

func sync_play_count_year(res []*Table) {
	arr := make([]*play_count_yearly_model.Table, 0, 400)
	for _, v := range res {
		st, sst := getShowType(v.Category)
		if st == -10 {
			continue
		}

		// 只记录全平台的数据
		if v.PlatformId != 0 {
			continue
		}
		row := &play_count_yearly_model.Table{
			ShowId:      v.ShowId,
			ShowType:    st,
			SubShowType: sst,
			Rank:        v.PlayCountRank,
			Rise:        v.PlayCountRise,
			Num:         v.PlayCount,
			DayAt:       model.StrToTime(v.Day),
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			play_count_yearly_model.Model().Create(&arr)
			arr = make([]*play_count_yearly_model.Table, 0, 400)
		}
	}
	if len(arr) > 0 {
		play_count_yearly_model.Model().Create(&arr)
	}
}

func getShowType(category string) (int64, int64) {
	switch category {
	case "ALL_ANIME":
		return show_model.ShowTypeAmine, -1
	case "ANIME":
		return show_model.ShowTypeAmine, show_model.ShowSubTypeAmineChina
	case "JAPAN_ANIME":
		return show_model.ShowTypeAmine, show_model.ShowSubTypeAmineJP
	case "DRAMA":
		return show_model.ShowTypeSeries, -1
	case "TV_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesTV
	case "NETWORK_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesNet
	case "ALL_MOVIE":
		return show_model.ShowTypeMovie, -1
	case "MOVIE":
		return show_model.ShowTypeMovie, show_model.ShowSubTypeMovieCinema
	case "NETWORK_MOVIE":
		return show_model.ShowTypeMovie, show_model.ShowSubTypeMovieNet
	case "VARIETY":
		return show_model.ShowTypeVariety, -1
	case "NETWORK_VARIETY":
		return show_model.ShowTypeVariety, show_model.ShowSubTypeVarietyNet
	case "TV_VARIETY":
		return show_model.ShowTypeVariety, show_model.ShowSubTypeVarietyTV
	}

	return -10, -10
}

func getType(tp string) int8 {
	switch tp {
	case "DAILY":
		return model_clean.CycleDaily
	case "WEEKLY":
		return model_clean.CycleWeekly
	case "MONTHLY":
		return model_clean.CycleMonthly
	}

	return -10
}

//func play_count_tatal ()  {
//	var res []*Table
//	Model().Select("*").Where("day >= ?",internal.StartDay).
//		Where("type", "TOTAL").
//		Find(&res)
//
//
//	arr := make([]*play_count_total_daily_model.Table, 0, 400)
//	for _, v := range res {
//		st, sst := getShowType(v.Category)
//		if st == -10 {
//			continue
//		}
//
//		// 只记录全平台的数据
//		if v.PlatformId != 0 {
//			continue
//		}
//		row := &play_count_total_daily_model.Table{
//			ShowId:   v.ShowId,
//			ShowType: st,
//			SubShowType: sst,
//			Rank: v.PlayCountRank,
//			Rise: v.PlayCountRise,
//			Num: v.PlayCount,
//			DayAt: model.StrToTime(v.Day),
//		}
//		arr = append(arr, row)
//		if len(arr) >= 400 {
//			play_count_total_daily_model.Model().Create(&arr)
//			arr = make([]*play_count_total_daily_model.Table, 0, 400)
//		}
//	}
//	if len(arr) > 0 {
//		play_count_total_daily_model.Model().Create(&arr)
//	}
//}
