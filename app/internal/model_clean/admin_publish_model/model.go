package admin_publish_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	IsShow model.Tinyint `json:"is_show"`
	Content model.Varchar `json:"content"`
	Type model.Varchar `json:"type"`
}

func (t Table) TableName() string {
	return "admin_publish"
}
