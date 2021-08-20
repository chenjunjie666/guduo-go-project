package play_count_total_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean/play_count_daily_model"
	"guduo/app/internal/model_clean/play_count_total_daily_model"
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

type Year struct {
	ShowId uint64
	Num int64
	NumMin int64
	DayAt uint
}

func Sync () {
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").
			Where("day >= ?",internal.StartDay).
			Where("type", "TOTAL").
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
	arr := make([]*play_count_total_daily_model.Table, 0, 400)
	arr2 := make([]*play_count_daily_model.Table, 0, 400)
	for _, v := range res {
		st, sst := getShowType(v.Category)
		if st == -10 {
			continue
		}

		// 只记录全平台的数据
		if v.PlatformId == 0 {
			row := &play_count_total_daily_model.Table{
				ShowId:   v.ShowId,
				ShowType: st,
				SubShowType: sst,
				Rank: v.TotalPlayCountRank,
				Rise: v.TotalPlayCountRise,
				Num: v.TotalPlayCount,
				DayAt: model.StrToTime(v.Day),
			}
			arr = append(arr, row)
			if len(arr) >= 400 {
				play_count_total_daily_model.Model().Create(&arr)
				arr = make([]*play_count_total_daily_model.Table, 0, 400)
			}
		}else if sst != -1 {
			row2 := &play_count_daily_model.Table{
				ShowId:   v.ShowId,
				PlatformId:  v.PlatformId,
				Num: v.TotalPlayCount,
				DayAt: model.StrToTime(v.Day),
			}
			arr2 = append(arr2, row2)
			if len(arr2) >= 400 {
				play_count_daily_model.Model().Create(&arr2)
				arr2 = make([]*play_count_daily_model.Table, 0, 400)
			}
		}
	}
	if len(arr) > 0 {
		play_count_total_daily_model.Model().Create(&arr)
	}
	if len(arr2) > 0 {
		play_count_daily_model.Model().Create(&arr2)
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