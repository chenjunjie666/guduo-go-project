package show

import (
	"guduo/app/internal/model_scrawler/show_model"
)

func List(showType, status int, keyword string, page, limit int) ([]map[string]interface{}, int64) {
	if page <= 0 {
		page = 1
	}
	if limit <= 10 {
		limit = 10
	}

	offset := (page - 1) * limit


	var res []*show_model.Table
	mdl := show_model.Model().Select(
		"id",
		"name",
		"show_type",
		"status",
		"release_at",
		"is_show",
		)

	mdl2 := show_model.Model()

	if keyword != "" {
		mdl = mdl.Where("name LIKE ?", "%" + keyword + "%")
		mdl2 = mdl2.Where("name LIKE ?", "%" + keyword + "%")
	}

	if showType >= 0 {
		mdl = mdl.Where("show_type", showType)
		mdl2 = mdl2.Where("show_type", showType)
	}

	if status != -99 {
		mdl = mdl.Where("status", status)
		mdl2 = mdl2.Where("status", status)
	}

	mdl = mdl.Offset(offset).
		Limit(limit).
		Order("id desc").
		Find(&res)

	var total int64
	mdl2.Count(&total)

	ret := make([]map[string]interface{}, len(res))

	for k, row := range res {
		isShowStr := "是"
		if row.IsShow != 1 {
			isShowStr = "否"
		}

		ret[k] = map[string]interface{}{
			"show_id": row.ID,
			"name": row.Name,
			"type_str": show_model.GetShowTypeStr(row.ShowType),
			"release_at": row.ReleaseAt,
			"show_status_str": show_model.GetShowStatusStr(row.ShowStatus),
			"is_show_str": isShowStr,
			"status_str": show_model.GetStatusStr(row.Status),
		}
	}

	return ret, total
}
