package guduo_hot_rank_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId      model.ForeignKey
	RankType    model.Tinyint
	ShowType    model.Int `json:"show_type"`
	SubShowType model.Int `json:"sub_show_type"`
	PlatformId  model.Int `json:"platform_id"`
	Num         model.Float
	CustomNum   model.Float
	Rank        model.Int
	Rise        model.Int
	DayAt       model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "guduo_hot_rank"
}
