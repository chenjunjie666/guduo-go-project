package show

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"guduo/app/cms/internal/hepler/validate"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_detail_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/errors"
)

type ShowDetailParams struct {
	ShowId       uint64 `json:"show_id"`       // 剧集ID
	ShowType     int64  `json:"show_type"`     // 一级分类
	SubShowType  int64  `json:"sub_show_type"` // 二级分类
	Name         string `json:"name"`          // 剧集名称
	Poster       string `json:"poster"`        // 海报url
	Introduction string `json:"introduction"`  // 剧集简介
	IsSelf       int8   `json:"is_self"`       // 是否自制
	IsAdapt      int8   `json:"is_adapt"`      // 是否改编
	AdaptFrom    string `json:"adapt_from"`    // 改编自
	ReleaseAt    uint   `json:"release_at"`    // 上线时间
	EndAt        uint   `json:"end_at"`        // 完结时间
	TotalEpisode int64  `json:"total_episode"` // 集数
	Length       string `json:"length"`        // 单集时长
	Staff        string `json:"staff"`         // 工作人员
	Director     string `json:"director"`      // 导演
	Actor        string `json:"actor"`         // 演员
	ShowStatus   int8   `json:"show_status"`   // 放映状态
	Link         string `json:"link"`          // 视频平台链接
	Tag          string `json:"tag"`           // 标签
	Status       int64  `json:"status"`        // 审核是否通过
	IsShow       int8   `json:"is_show"`       // 是否在前端显示
}

func Detail(sid uint64) interface{} {
	var res show_model.Table

	show_model.Model().Where("id", sid).Find(&res)

	var actors []show_actor_model.Table
	show_actor_model.Model().Where("show_id", sid).Find(&actors)

	var urls []show_detail_model.Table
	show_detail_model.Model().Where("show_id", sid).Find(&urls)

	ret := map[string]interface{}{
		"show_id":       res.ID,
		"name":          res.Name,
		"poster":        res.Poster,
		"show_type":     res.ShowType,
		"sub_show_type": res.SubShowType,
		"platform":      show_model.GetPlatform(res.Platform),
		"status":        res.Status,
		"tag":           show_model.GetTag(res.Tag),
		"introduction":  res.Introduction,
		"staff":         show_model.GetStaffWithoutDirector(res.Staff),
		"director":      show_model.GetDirector(res.Director),
		"length":        res.Length,
		"release_at":    res.ReleaseAt,
		"end_at":        res.EndAt,
		"total_episode": res.TotalEpisode,
		"is_show":       res.IsShow,
		"is_self":       res.IsSelf,
		"is_adapt":      res.IsAdapt,
		"adapt_from":    res.AdaptFrom,
		"show_status":   res.ShowStatus,
		"link":          urls,
		"actor":         actors,
	}
	return ret
}

func AddNewShow(data ShowDetailParams) (uint64, error) {
	var detailLink []*show_detail_model.Table
	e := json.Unmarshal([]byte(data.Link), &detailLink)

	if e != nil {
		return 0, errors.CmsError("详情链接不正确")
	}
	allPlt := constant.GetVideoPlatformMap()
	for _, link := range detailLink {
		b := validate.ShowDetailUrl(link.PlatformId, link.Url)
		if !b {
			return 0, errors.CmsError(fmt.Sprintf("平台:%s 的链接不正确", allPlt[int(link.PlatformId)]))
		}
	}

	// 播放的平台ID
	plt := make([]uint64, 0, 2)
	for _, link := range detailLink {
		for pid := range allPlt {
			lpid := link.PlatformId
			if lpid == uint64(pid) {
				plt = append(plt, lpid)
			}
		}
	}
	pltJson, _ := json.Marshal(plt)

	if data.Tag == "" {
		data.Tag = "[]"
	}

	isIntro := int8(0)
	if data.Introduction != "" {
		isIntro = 1
	}

	isLen := int8(0)
	if data.Length != "" {
		isLen = 1
	}

	isBaseInfo := int8(9)
	if data.Staff != "" || data.Director != "" {
		isBaseInfo = 1
	}

	isRelase := int8(0)
	if data.ReleaseAt > 0 {
		isRelase = 1
	}

	var tag show_model.Tag
	e = json.Unmarshal([]byte(data.Staff), &tag)
	if e != nil {
		tag = make(show_model.Tag, 0)
	}
	tagJson, _ := json.Marshal(tag)

	var staff show_model.StaffWithoutDirector
	e = json.Unmarshal([]byte(data.Staff), &staff)
	if e != nil {
		staff = show_model.EmptyStaffWithoutDirector()
	}
	staffJson, _ := json.Marshal(staff)

	var director show_model.Director
	e = json.Unmarshal([]byte(data.Director), &director)
	if e != nil {
		director = make(show_model.Director, 0)
	}
	directorJson, _ := json.Marshal(director)

	var actors []*show_actor_model.Table
	e = json.Unmarshal([]byte(data.Actor), &actors)
	if e != nil {
		actors = make([]*show_actor_model.Table, 0)
	}

	save := show_model.Table{
		Name:              data.Name,
		Poster:            data.Poster,
		ShowType:          data.ShowType,
		SubShowType:       data.SubShowType,
		Platform:          string(pltJson),
		Status:            data.Status,
		Tag:               string(tagJson),
		Introduction:      data.Introduction,
		Staff:             string(staffJson),
		Director:          string(directorJson),
		Length:            data.Length,
		ReleaseAt:         data.ReleaseAt,
		EndAt:             data.EndAt,
		TotalEpisode:      data.TotalEpisode,
		IsCrawlerBaseInfo: isBaseInfo,
		IsCrawlerIntro:    isIntro,
		IsCrawlerLen:      isLen,
		IsCrawlerRelease:  isRelase,
		IsShow:            data.IsShow,
		IsSelf:            data.IsSelf,
		IsAdapt:           data.IsAdapt,
		AdaptFrom:         data.AdaptFrom,
		ShowStatus:        data.ShowStatus,
	}
	r := show_model.Model().Create(&save)

	if r.Error != nil {
		log.Warn(r.Error)
		return 0, errors.CmsError("新增剧集失败")
	}

	updateLinks(detailLink, save.ID)
	updateActors(actors, save.ID)

	return save.ID, nil
}

func Update(data ShowDetailParams) error {
	sid := data.ShowId

	var detailLink []*show_detail_model.Table
	e := json.Unmarshal([]byte(data.Link), &detailLink)

	if e != nil {
		return errors.CmsError("详情链接不正确")
	}

	allPlt := constant.GetVideoPlatformMap()
	for _, link := range detailLink {
		b := validate.ShowDetailUrl(link.PlatformId, link.Url)
		if !b {
			return errors.CmsError(fmt.Sprintf("平台:%s 的链接不正确", allPlt[int(link.PlatformId)]))
		}
	}

	// 播放的平台ID
	plt := make([]uint64, 0, 2)
	for _, link := range detailLink {
		for pid := range allPlt {
			lpid := link.PlatformId
			if lpid == uint64(pid) {
				plt = append(plt, lpid)
			}
		}
	}
	pltJson, _ := json.Marshal(plt)

	var tag show_model.Tag
	e = json.Unmarshal([]byte(data.Tag), &tag)
	if e != nil {
		tag = make(show_model.Tag, 0)
	}
	tagJson, _ := json.Marshal(tag)

	var staff show_model.StaffWithoutDirector
	e = json.Unmarshal([]byte(data.Staff), &staff)
	if e != nil {
		staff = show_model.EmptyStaffWithoutDirector()
	}
	staffJson, _ := json.Marshal(staff)

	var director show_model.Director
	e = json.Unmarshal([]byte(data.Director), &director)
	if e != nil {
		director = make(show_model.Director, 0)
	}
	directorJson, _ := json.Marshal(director)

	var actors []*show_actor_model.Table
	e = json.Unmarshal([]byte(data.Actor), &actors)
	if e != nil {
		actors = make([]*show_actor_model.Table, 0)
	}

	update := show_model.Table{
		Name:         data.Name,
		Poster:       data.Poster,
		ShowType:     data.ShowType,
		SubShowType:  data.SubShowType,
		Platform:     string(pltJson),
		Status:       data.Status,
		Tag:          string(tagJson),
		Introduction: data.Introduction,
		Staff:        string(staffJson),
		Director:     string(directorJson),
		Length:       data.Length,
		ReleaseAt:    data.ReleaseAt,
		EndAt:        data.EndAt,
		TotalEpisode: data.TotalEpisode,
		IsShow:       data.IsShow,
		IsSelf:       data.IsSelf,
		IsAdapt:      data.IsAdapt,
		AdaptFrom:    data.AdaptFrom,
		ShowStatus:   data.ShowStatus,
	}

	if data.Introduction != "" {
		update.IsCrawlerIntro = 1
	}
	if data.Length != "" {
		update.IsCrawlerLen = 1
	}
	if data.Staff != "" || data.Director != "" {
		update.IsCrawlerBaseInfo = 1
	}
	if data.ReleaseAt > 0 {
		update.IsCrawlerRelease = 1
	}

	update.ID = sid
	r := show_model.Model().Select("*").
		Omit("is_crawler_base_info", "is_crawler_intro", "is_crawler_len", "is_crawler_release").
		Where("id", sid).Updates(&update)

	if r.Error != nil {
		log.Warn(r.Error)
		return errors.CmsError("更新信息失败")
	}
	updateLinks(detailLink, sid)
	updateActors(actors, sid)
	return nil
}

func updateLinks(detailLink []*show_detail_model.Table, sid uint64) {
	for k := range detailLink {
		detailLink[k].ShowId = sid
		detailLink[k].Usable = 1
	}

	show_detail_model.Model().Where("show_id", sid).Delete(nil)

	if len(detailLink) == 0 {
		return
	}
	show_detail_model.Model().Create(&detailLink)
}

func updateActors(showActorsUpdate []*show_actor_model.Table, sid uint64) {
	if len(showActorsUpdate) == 0 {
		return
	}

	names := make([]string, len(showActorsUpdate))
	for _, row := range showActorsUpdate {
		names = append(names, row.Name)
	}

	var actors []actor_model.Table
	actor_model.Model().Select("id", "name").Where("name IN ?", names).
		Find(&actors)

	for k, v := range showActorsUpdate {
		showActorsUpdate[k].ShowId = sid
		for _, actor := range actors {
			if actor.Name == v.Name {
				showActorsUpdate[k].ActorId = v.ID
				break
			}
		}
	}
	show_actor_model.Model().Where("show_id", sid).Delete(nil)
	if len(showActorsUpdate) == 0 {
		return
	}

	show_actor_model.Model().Create(&showActorsUpdate)
}
