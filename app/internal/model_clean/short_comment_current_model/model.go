package article_content_current_model

import "guduo/pkg/model"

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	PlatformId model.ForeignKey
	Content model.Text
	Author model.Varchar
	PublishAt model.SecondTimeStamp
	Forward model.Int
	JobAt model.SecondTimeStamp
}

func (t Table) TableName() string {
	return "article_content_current"
}
