package show

import (
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/time"
)

type baseInfo struct {
	ShowId         uint64                      `json:"show_id"`
	Name           string                      `json:"name"`
	ShowType       int64                       `json:"show_type"`
	ShowTypeStr    string                      `json:"show_type_str"`
	SubShowType    int64                       `json:"sub_show_type"`
	SubShowTypeStr string                      `json:"sub_show_type_str"`
	Poster         string                      `json:"poster"`
	Platform       show_model.Platform         `json:"platform"`
	Staff          show_model.Staff            `json:"staff"`
	Tag            show_model.Tag              `json:"tag"`
	Introduction   string                      `json:"introduction"`
	Length         string                      `json:"length"`
	ReleaseAt      string                      `json:"release_at"`
	Released       int64                       `json:"released"`
	TotalEpisode   int64                       `json:"total_episode"`
	Actor          show_actor_model.ShowActors `json:"actor"`
	ShowStatus     int8                        `json:"show_status"`
}

func GetBaseInfo(sid uint64) *baseInfo {
	bInfo := show_model.GetBaseInfo(sid)
	actors := show_actor_model.GetShowActors(sid)

	ret := &baseInfo{
		ShowId:         bInfo.ID,
		Name:           bInfo.Name,
		ShowType:       bInfo.ShowType,
		ShowTypeStr:    show_model.GetShowTypeStr(bInfo.ShowType),
		SubShowType:    bInfo.SubShowType,
		SubShowTypeStr: show_model.GetSubShowTypeStr(bInfo.SubShowType),
		Poster:         bInfo.Poster,
		Platform:       show_model.GetPlatform(bInfo.Platform),
		Staff:          show_model.GetStaff(bInfo.Staff, bInfo.Director),
		Tag:            show_model.GetTag(bInfo.Tag),
		Introduction:   bInfo.Introduction,
		Length:         bInfo.Length,
		ReleaseAt:      time.TimeToStr(time.LayoutYmd, bInfo.ReleaseAt),
		Released:       int64((time.Today() - bInfo.ReleaseAt) / 86400),
		TotalEpisode:   bInfo.TotalEpisode,
		ShowStatus:     bInfo.ShowStatus,
		Actor:          actors,
	}

	return ret
}
