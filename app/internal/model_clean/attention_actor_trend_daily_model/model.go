package attention_actor_trend_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ActorId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Int
	CustomNum model.Int
	Post model.Int
	CustomPost model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "attention_actor_trend_daily"
}
