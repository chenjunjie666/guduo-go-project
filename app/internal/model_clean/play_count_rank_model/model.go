package play_count_rank_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	ShowType model.Int
	SubShowType model.Int
	RankType model.Tinyint
	Num model.Int
	Rank model.Int
	Rise model.Int
	DayAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "play_count_rank"
}
