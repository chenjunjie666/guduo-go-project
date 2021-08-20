package play_count_daily_inc_model
//
//import (
//	"gorm.io/gorm"
//	play_count_trend_daily_model "guduo/app/internal/model_clean/play_count_trend_daily_model"
//	"guduo/app/migration/internal"
//	"guduo/app/migration/model"
//	"guduo/pkg/db"
//)
//
//var m *gorm.DB
//
//func Model() *gorm.DB {
//	if m == nil {
//		m = db.GetLoliPopMysqlConn()
//	}
//	return m.Model(&Table{})
//}
//
//func Sync () {
//	var res []*Table
//	Model().Select("*").Where("day >= ?",internal.StartDay).Find(&res)
//
//	arr := make([]*play_count_trend_daily_model.Table, 0, 400)
//	for _, v := range res {
//		row := &play_count_trend_daily_model.Table{
//			ShowId:   v.ShowId,
//			PlatformId:  v.PlatformId,
//			Num: v.PlayCount,
//			DayAt: model.StrToTime(v.Day),
//		}
//		arr = append(arr, row)
//		if len(arr) >= 400 {
//			play_count_trend_daily_model.Model().Create(&arr)
//			arr = make([]*play_count_trend_daily_model.Table, 0, 400)
//		}
//	}
//	if len(arr) > 0 {
//		play_count_trend_daily_model.Model().Create(&arr)
//	}
//}
