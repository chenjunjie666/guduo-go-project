package show_actor_model

import "guduo/pkg/model"

type Table struct {
	model.Id
	ShowId   model.ForeignKey `json:"show_id"`
	ActorId  model.ForeignKey `json:"actor_id"`
	Name     model.Varchar `json:"name"`
	Avatar   model.Varchar `json:"avatar"`
	Play     model.Varchar `json:"play"`
	PlayType model.Tinyint `json:"play_type"`
}

func (t Table) TableName() string {
	return "show_actor"
}
