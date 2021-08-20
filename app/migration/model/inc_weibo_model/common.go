package inc_weibo_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/article_count_trend_daily_model"
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
	offset := 0
	limit := 100000
	for {
		var res []*Table
		Model().Select("*").Where("day >= ?", internal.StartDay).
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
	arr := make([]*article_count_trend_daily_model.Table, 0, 400)
	for _, v := range res {
		row := &article_count_trend_daily_model.Table{
			ShowId:     v.ShowId,
			PlatformId: constant.PlatformIdWeibo,
			Num:        v.Count,
			DayAt:      model.StrToTime(v.Day),
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			article_count_trend_daily_model.Model().Create(&arr)
			arr = make([]*article_count_trend_daily_model.Table, 0, 400)
		}
	}
	if len(arr) > 0 {
		article_count_trend_daily_model.Model().Create(&arr)
	}
}
