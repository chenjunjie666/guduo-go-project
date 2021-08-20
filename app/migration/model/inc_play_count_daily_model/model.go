package inc_play_count_daily_model

import "guduo/pkg/model"

type Table struct {
	Name          string
	ShowId        model.ForeignKey
	PlatformId    model.ForeignKey
	Category      model.Varchar
	PlayCount     model.Int
	PlayCountRank model.Int
	PlayCountRise model.Int
	Day           model.Varchar
}

func (t Table) TableName() string {
	return "inc_billboard_combine_day_type_play_count"
	//return "inc_billboard_combine_day_type_play_count"
}
