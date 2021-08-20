package attention_actor_model

import (
	"guduo/pkg/model"
)

type Table struct {
	model.Fields
	PlatformId model.ForeignKey
	ActorId model.ForeignKey
	Attention model.Int
	Post model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "attention_actor"
}