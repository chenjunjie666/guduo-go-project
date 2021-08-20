package admin_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	Username model.Varchar `json:"username"`
	Password model.Varchar `json:"password"`
}

func (t Table) TableName() string {
	return "admin"
}
