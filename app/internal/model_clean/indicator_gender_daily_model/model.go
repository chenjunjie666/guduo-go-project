package indicator_gender_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId       model.ForeignKey
	PlatformId   model.ForeignKey
	MaleRating   model.Float
	FemaleRating model.Float
	DayAt        model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "indicator_gender_daily"
}
