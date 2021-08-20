package tag_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	Name model.Varchar `json:"name"`
}

func (t Table) TableName() string {
	return "tag"
}
