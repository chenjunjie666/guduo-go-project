package indicator_age_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Rating model.Float
	CustomRating model.Float
	AgeFrom model.Int
	AgeTo model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "indicator_age_daily"
}
