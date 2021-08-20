package inc_weibo_model

import "guduo/pkg/model"

type Table struct {
	Name     string
	Category model.Varchar
	ShowId   model.ForeignKey
	Count    model.Int
	Day      model.Varchar
}

func (t Table) TableName() string {
	return "inc_billboard_ranking_weibo"
}
