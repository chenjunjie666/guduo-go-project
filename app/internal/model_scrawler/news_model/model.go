package news_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "news"
}
