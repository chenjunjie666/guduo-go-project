package show_gdi_logs_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/guduo_hot_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_rank_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/app/migration/internal"
	model2 "guduo/app/migration/model"
	"guduo/pkg/db"
	"guduo/pkg/model"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetLoliPopMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync () {
	guduo_hot_rank_model.Model().Where("day_at >= ?", model2.StrToTime(internal.StartDay)).Delete(nil)
	guduo_hot_daily_model.Model().Where("day_at >= ?",  model2.StrToTime(internal.StartDay)).Delete(nil)
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").
			Where("day >= ?", internal.StartDay).
			Offset(offset).
			Limit(limit).
			Find(&res)

		sync(res)

		if len(res) < 100000 {
			break
		}

		offset += limit
	}


}


func sync(res []*Table) {

	arr := make([]*guduo_hot_rank_model.Table, 0, 400)
	arr2 := make([]*guduo_hot_daily_model.Table, 0, 400)
	for _, v := range res {
		type1, type2 :=  getShowType(v.Category)

		if type1 == -10 {
			continue
		}

		rank_type := getType(v.Type)
		if rank_type == -10 {
			continue
		}

		if rank_type == model_clean.CycleDaily {
			row2 := &guduo_hot_daily_model.Table{
				ShowId:    v.ShowId,
				Num:       v.Gdi,
				DayAt:    model2.StrToTime(v.Day),
			}
			arr2 = append(arr2, row2)
			if len(arr2) >= 400 {
				guduo_hot_daily_model.Model().Create(&arr2)
				arr2 = arr2[:0]
			}
		}


		row := &guduo_hot_rank_model.Table{
			ShowId:   v.ShowId,
			ShowType:  type1,
			SubShowType:  type2,
			PlatformId:  v.PlatformId,
			Num: v.Gdi,
			Rank: v.GdiRank,
			RankType: rank_type,
			Rise: v.GdiRise,
			DayAt: model2.StrToTime(v.Day),
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			guduo_hot_rank_model.Model().Create(&arr)
			arr = arr[:0]
		}
	}

	if len(arr) > 0 {
		guduo_hot_rank_model.Model().Create(&arr)
	}

	if len(arr2) >= 400 {
		guduo_hot_daily_model.Model().Create(&arr2)
		arr2 = arr2[:0]
	}
}


func getShowType(category model.Varchar) (int64, int64) {
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