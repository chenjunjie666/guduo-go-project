package article_content_model

import (
	"guduo/pkg/model"
)

type Table struct {
	model.Fields
	PlatformId model.ForeignKey
	ShowId model.ForeignKey
	UID model.Varchar
	Content model.Text
	Author model.Varchar
	PublishAt model.SecondTimeStamp
	Forward model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "article_content"
}


