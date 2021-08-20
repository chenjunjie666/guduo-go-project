package article_count_trend_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "article_count_trend_daily"
}
