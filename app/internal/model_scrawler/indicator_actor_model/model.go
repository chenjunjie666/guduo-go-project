package indicator_actor_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ActorId model.ForeignKey
	PlatformId model.ForeignKey
	Num model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "indicator_actor"
}
