package actor_billboard_model


import "guduo/pkg/model"

type Table struct {
	Name string
	ActorId   model.ForeignKey // ymd -> end at
	ActorName   model.Varchar // ymd -> release at
	EffectionIndex      model.Float // string -> show type
	Category     model.Varchar // a,b,c -> platform
	BillboardRank      model.Int // 时长-分钟 -> string length
	EffectionRankRise model.Int
	Day model.Varchar
}

func (t Table) TableName() string {
	return "actor_billboard"
}