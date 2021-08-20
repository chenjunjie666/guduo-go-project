package article_count_actor_daily_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ActorId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "article_count_actor_daily"
}
