package douban_logs_model

import "guduo/pkg/model"

type Table struct {
	Name      string
	ShowId    model.ForeignKey
	Score    model.Float
	Day model.Varchar
}

func (t Table) TableName() string {
	return "douban_logs"
}