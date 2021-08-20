package carl_shows_model

import "guduo/pkg/model"

type Table struct {
	PlatformId    model.ForeignKey
	LinkedId model.ForeignKey
	Url model.Varchar
}

func (t Table) TableName() string {
	return "shows"
}