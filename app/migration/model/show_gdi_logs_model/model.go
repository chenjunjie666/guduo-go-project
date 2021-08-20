package show_gdi_logs_model

import "guduo/pkg/model"

type Table struct {
	Name      string
	ShowId    model.ForeignKey
	Type model.Varchar
	Category model.Varchar
	PlatformId model.Int
	Gdi   model.Float
	GdiRank model.Int
	GdiRise model.Int
	Day model.Varchar
}

func (t Table) TableName() string {
	return "inc_billboard_combine_base"
}
