package platform_model

import (
	"guduo/pkg/model"
)

type Table struct {
	model.Fields
	Url  model.Text
	Name model.Varchar
}

func (m Table) TableName() string {
	return "platform"
}
