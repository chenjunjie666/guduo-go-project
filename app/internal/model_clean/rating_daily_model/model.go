package rating_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Rating model.Float
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "rating_daily"
}
