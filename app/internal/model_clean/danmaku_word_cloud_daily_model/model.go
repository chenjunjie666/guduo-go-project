package danmaku_word_cloud_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	Word model.Varchar
	Weight model.Int
	CloudPic model.Text
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "danmaku_word_cloud_daily"
}
