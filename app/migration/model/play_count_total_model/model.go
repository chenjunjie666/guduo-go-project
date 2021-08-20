package play_count_total_model

import "guduo/pkg/model"

type Table struct {
	Name               string
	ShowId             model.ForeignKey
	PlatformId         model.ForeignKey
	Category           model.Varchar
	TotalPlayCount     model.Int
	TotalPlayCountRank model.Int
	TotalPlayCountRise model.Int
	Day                model.Varchar
}

func (t Table) TableName() string {
	return "inc_billboard_combine_total_play_count"
}
