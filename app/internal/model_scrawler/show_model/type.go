package show_model

import "guduo/pkg/model"

type ShowBaseInfo struct {
	ID model.PrimaryKey
	Name model.Varchar
	ShowType model.Int
	SubShowType model.Int
	Poster model.Varchar
	Platform model.Varchar
	Staff model.Text
	Tag model.Varchar
	Introduction model.Text
	Director model.Text
	Length model.Varchar
	ReleaseAt model.SecondTimeStamp
	TotalEpisode model.Int
	Status model.Int
	ShowStatus model.Tinyint
}


type ShowName struct {
	Id model.PrimaryKey
	Name model.Varchar
}

type Platform []int
type Tag []string
// 导演是通过爬虫爬取的，如果存在staff会对这个字段频繁读写并且进行
// 非常慢的json反序列化，所以在存储的时候和staff分开，在小程序，后台等地方
// 在做合并处理
type Director = []string

// 剧综的工作组人员
type Staff struct {
	Director Director `json:"director"`
	ScreenWriter []string `json:"screen_writer"` // 编剧 json数组
	Producer []string `json:"producer"` // 制片人 json数组
	ProducerCompany []string `json:"producer_company"` // 制片人公司 json数组
	Publisher []string `json:"publisher"` // 出品人 json数组
	PublisherCompany []string `json:"publisher_company"` // 出品公司 json 数组
}

// 剧综的工作组人员
type StaffWithoutDirector struct {
	ScreenWriter []string `json:"screen_writer"` // 编剧 json数组
	Producer []string `json:"producer"` // 制片人 json数组
	ProducerCompany []string `json:"producer_company"` // 制片人公司 json数组
	Publisher []string `json:"publisher"` // 出品人 json数组
	PublisherCompany []string `json:"publisher_company"` // 出品公司 json 数组
}