package word_cloud_model

import "guduo/pkg/model"

type Table struct {
	Name       string
	ShowId     model.ForeignKey
	PlatformId model.ForeignKey
	Weight     model.Int
	Content    model.Varchar
	CreateTime model.Varchar
}

func (t Table) TableName() string {
	return "barrage_wordcloud_logs"
}

