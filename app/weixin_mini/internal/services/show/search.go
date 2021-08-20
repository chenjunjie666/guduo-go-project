package show

import (
	"guduo/app/internal/model_scrawler/show_actor_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"strings"
)

func SearchShow(k string) []map[string]interface{} {
	type search struct {
		ID             model.PrimaryKey      `json:"id"`
		Name           model.Varchar         `json:"name"`
		Poster         model.Varchar         `json:"poster"`
		SubShowType    model.Int             `json:"sub_show_type"`
		Platform       model.Varchar         `json:"platform"`
		Director       model.Text            `json:"director"`
		ReleaseAt      model.SecondTimeStamp `json:"release_at"`
		ShowStatus     model.Tinyint         `json:"show_status"`
		Actor          model.Varchar         `json:"actor"`
	}

	var list []search
	show_model.Model().Select("id", "name", "poster", "release_at", "show_status", "director", "platform", "sub_show_type").
		Where("name LIKE ?", "%" + k +"%").
		Where("status", show_model.ShowStatStandard).
		Where("is_show", show_model.ShowOn).
		Order("id desc").
		Limit(50).
		Find(&list)

	sid := make([]uint64, len(list))
	for key, v := range list {
		sid[key] = v.ID
	}

	var actor []show_actor_model.Table
	show_actor_model.Model().Select("show_id", "GROUP_CONCAT(name SEPARATOR '/') as name").
		Where("show_id IN ?", sid).
		Group("show_id").
		Find(&actor)

	ret := make([]map[string]interface{}, len(list))
	for key, v := range list {
		name := ""
		for _, row := range actor {
			if row.ShowId == v.ID {
				name = row.Name
				actArr := strings.Split(name, "/")
				if len(actArr) > 5 {
					actArr = actArr[0:5]
					name = strings.Join(actArr, "/")
				}
			}
		}

		director := show_model.GetDirector(v.Director)

		dirStr := ""
		for _, d := range director {
			dirStr += d + "/"
		}
		dirStr = strings.Trim(dirStr, "/")

		ret[key] = map[string]interface{}{
			"show_id": v.ID,
			"name": v.Name,
			"poster": v.Poster,
			"sub_show_type": v.SubShowType,
			"sub_show_type_str": show_model.GetSubShowTypeStr(v.SubShowType),
			"platform": show_model.GetPlatform(v.Platform),
			"director": dirStr,
			"release_at": v.ReleaseAt,
			"actor": name,
			"show_status_str": show_model.GetShowStatusStr(v.ShowStatus),
		}
	}

	return ret
}


type search struct {
	ID             model.PrimaryKey      `json:"show_id"`
	Name           model.Varchar         `json:"name"`
}

func HotSearch(type_ int) []search {
	var list []search
	show_model.Model().Select("id", "name").
		Where("is_search_hot", show_model.ShowSearchHot).
		Where("status", show_model.ShowStatStandard).
		Where("is_show", show_model.ShowOn).
		Where("show_type", type_).
		Limit(15).
		Find(&list)

	return list
}