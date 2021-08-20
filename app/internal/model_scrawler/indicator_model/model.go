package indicator_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Float
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "indicator"
}
