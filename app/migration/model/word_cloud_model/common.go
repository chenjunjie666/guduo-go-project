package word_cloud_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean/danmaku_word_cloud_daily_model"
	"guduo/app/migration/internal"
	model2 "guduo/app/migration/model"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetLoliPopMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync () {
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").
			Where("create_time > ?", internal.StartTime).
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

func sync(res []*Table){

	arr := make([]*danmaku_word_cloud_daily_model.Table, 0, 400)
	for _, v := range res {
		row := &danmaku_word_cloud_daily_model.Table{
			ShowId:   v.ShowId,
			Word: v.Content,
			Weight: v.Weight,
			DayAt: model2.StrToTime(v.CreateTime),
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			danmaku_word_cloud_daily_model.Model().Create(&arr)
			arr = make([]*danmaku_word_cloud_daily_model.Table, 0, 400)
		}
	}
	if len(arr) > 0 {
		danmaku_word_cloud_daily_model.Model().Create(&arr)
	}
}
