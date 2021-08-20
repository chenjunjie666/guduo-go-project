package danmaku_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Content model.Varchar
	ContentId model.Varchar
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "danmaku"
}
