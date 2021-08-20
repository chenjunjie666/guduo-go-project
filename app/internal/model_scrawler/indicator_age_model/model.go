package indicator_age_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	AgeFrom model.Int
	AgeTo model.Int
	Rating model.Float
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "indicator_age"
}
