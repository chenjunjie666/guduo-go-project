package play_count_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	PlatformId model.ForeignKey
	ShowId model.ForeignKey
	Num model.Int
	JobAt model.SecondTimeStamp
}

func (d Table) TableName() string {
	return "play_count"
}
