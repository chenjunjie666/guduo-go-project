package show_model

import (
	"guduo/pkg/model"
)

type Table struct {
	model.FieldsWithSoftDelete
	Name              model.Varchar         `json:"name"`
	Poster            model.Varchar         `json:"poster"`
	ShowType          model.Int             `json:"show_type"`
	SubShowType       model.Int             `json:"sub_show_type"`
	Platform          model.Varchar         `json:"platform"`
	Status            model.Int             `json:"status"`
	Tag               model.Varchar         `json:"tag"`
	Introduction      model.Text            `json:"introduction"`
	Staff             model.Varchar         `json:"staff"`
	Director          model.Text            `json:"director"`
	Length            model.Varchar         `json:"length"`
	ReleaseAt         model.SecondTimeStamp `json:"release_at"`
	EndAt             model.SecondTimeStamp `json:"end_at"`
	TotalEpisode      model.Int             `json:"total_episode"`
	IsCrawlerBaseInfo model.Tinyint         `json:"is_crawler_base_info"`
	IsCrawlerIntro    model.Tinyint         `json:"is_crawler_intro"`
	IsCrawlerLen      model.Tinyint         `json:"is_crawler_len"`
	IsCrawlerRelease  model.Tinyint         `json:"is_crawler_release"`
	IsShow            model.Tinyint         `json:"is_show"`
	IsSelf            model.Tinyint         `json:"is_self"`
	IsAdapt           model.Tinyint         `json:"is_adapt"`
	AdaptFrom         model.Varchar         `json:"adapt_from"`
	ShowStatus        model.Tinyint         `json:"show_status"`
}

func (d Table) TableName() string {
	return "show"
}
