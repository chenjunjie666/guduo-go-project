package actor_hot_monthly_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ActorId model.ForeignKey
	Num model.Float
	CustomNum model.Float
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "actor_hot_monthly"
}
