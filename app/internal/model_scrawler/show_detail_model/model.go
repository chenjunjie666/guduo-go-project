package show_detail_model

import (
	"guduo/pkg/model"
)

type Table struct {
	model.Fields
	ShowId model.ForeignKey `json:"show_id"`
	PlatformId model.ForeignKey `json:"platform_id"`
	Url model.Varchar `json:"url"`
	Usable model.Int `json:"usable"`
	TrueUrl model.Varchar `json:"true_url"`
}

func (d Table) TableName() string {
	return "show_detail"
}
