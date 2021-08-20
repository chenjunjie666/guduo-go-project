package bi_shows_model

import "guduo/pkg/model"

type Table struct {
	ID               model.PrimaryKey `json:"id" gorm:"primaryKey"`
	Name             string
	ReleaseDate      model.Varchar // ymd -> release at
	OfflineDate      model.Varchar // ymd -> end at
	Category         model.Varchar // string -> show type
	Platforms        model.Varchar // a,b,c -> platform
	Duration         model.Varchar // 时长-分钟 -> string length
	Episode          model.Int     // total_episode
	CoverImgUrl      model.Varchar //poster
	Director         model.Varchar // director
	Type             model.Varchar // director
	Intro            model.Varchar // introduction
	Locked           model.BitBool // is show
	ReleaseStatus    model.Tinyint // show_status
	MadeInSelf       model.Tinyint // is_self
	PremiereDate     model.Varchar // 首映日
	Publisher        model.Varchar // a,b,c
	Producer         model.Varchar // a,b,c
	ScriptWriter     model.Varchar // a,b,c
	AdaptedStatus    model.Tinyint // is adapt
	AdaptedFrom      model.Varchar // 推测是改编类型，暂无字段
	AdaptedWorksName model.Varchar // adapt from
}

func (t Table) TableName() string {
	return "shows"
}
