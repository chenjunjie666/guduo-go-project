package actor_domi_rank_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ActorId model.ForeignKey
	ActorName model.Varchar
	PlayType model.Tinyint
	Cycle model.Tinyint
	Num model.Float
	CustomNum model.Float
	Rank model.Int
	Rise model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "actor_domi_rank"
}
