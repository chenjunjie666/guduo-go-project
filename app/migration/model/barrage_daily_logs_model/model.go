package barrage_daily_logs_model

import "guduo/pkg/model"

type Table struct {
	Name      string
	ShowId    model.ForeignKey
	PlatformId   model.ForeignKey
	Count    model.Int
	Day model.Varchar
}

func (t Table) TableName() string {
	return "barrage_daily_logs"
}