package indicator_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Float
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "indicator_daily"
}
