package article_num_model

import (
	"guduo/pkg/model"
)

type Table struct {
	model.Fields
	PlatformId model.ForeignKey
	ShowId model.ForeignKey
	Num model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "article_num"
}