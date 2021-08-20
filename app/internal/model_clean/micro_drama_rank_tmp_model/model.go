package micro_drama_rank_tmp_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	//ShowId model.ForeignKey
	Name model.Varchar `json:"name"`
	PlatformId model.Varchar `json:"platform_id"`
	Num model.Float `json:"num"`
	DayAt model.SecondTimeStamp `json:"day_at"`
}

func (t Table) TableName() string {
	return "micro_drama_rank_tmp"
}
