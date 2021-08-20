package r_show_actor_model


import "guduo/pkg/model"

type Table struct {
	ID model.PrimaryKey `json:"id" gorm:"primaryKey"`
	Name string
	ShowId   model.ForeignKey // ymd -> release at
	ActorId   model.ForeignKey // ymd -> end at
	Avatar      model.Varchar // string -> show type
	ActorName     model.Varchar // a,b,c -> platform
	Roles      model.Varchar // 时长-分钟 -> string length
	Category model.Varchar
}

func (t Table) TableName() string {
	return "r_show_actor"
}