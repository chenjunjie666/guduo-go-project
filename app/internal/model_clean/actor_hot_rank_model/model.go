package actor_hot_rank_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ActorId model.ForeignKey
	ActorName model.Varchar
	IsNew model.Tinyint
	Cycle model.Tinyint
	Num model.Float
	CustomNum model.Float
	Rank model.Int
	Rise model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "actor_hot_rank"
}
