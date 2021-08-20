package guduo_hot_avg_rank_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	Num model.Float
	CustomNum model.Float
	Rank model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "guduo_hot_avg_rank"
}
