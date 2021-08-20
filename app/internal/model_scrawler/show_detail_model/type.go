package show_detail_model

import "guduo/pkg/model"

type DetailLink struct {
	ShowId model.ForeignKey `json:"show_id"`
	PlatformId model.ForeignKey `json:"platform_id"`
	Url string `json:"url"`
}

type DetailLinks []DetailLink