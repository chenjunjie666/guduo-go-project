package hot_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Hot model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "hot"
}
