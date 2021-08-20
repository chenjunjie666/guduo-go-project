package show

import (
	"fmt"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/guduo_hot_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_rank_model"
	"guduo/app/internal/model_clean/play_count_rank_model"
	"guduo/app/internal/model_clean/play_count_total_daily_model"
	"guduo/app/internal/model_clean/play_count_yearly_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"guduo/pkg/time"
)

const (
	ListTypeDay = iota
	ListTypeWeek
	ListTypeMonth
	ListTypeYear
	ListTypeAll
)

type HotTable struct {
	model.Fields
	ShowId    model.ForeignKey
	Num       model.Float
	CustomNum model.Float
	DayAt     model.SecondTimeStamp
	Days      model.Int
	Trend     model.Int
}

type PlayCountTable struct {
	model.Fields
	ShowId     model.ForeignKey
	PlatformId model.ForeignKey
	Num        model.Int
	CustomNum  model.Int
	DayAt      model.SecondTimeStamp
}

func List(dayAt, listType, type_, subType, pid int) []map[string]interface{} {
	// 先找到对应条件下的 show id，榜单数据只在这些 show id 下找
	ret := guduoRank(dayAt, listType, type_, subType, pid)
	return ret
}

// 获取剧集ID
func showIds(type_, subType, pid int) []uint64 {
	var res []show_model.Table
	mdl := show_model.Model()
	mdl = mdl.Select("id")
	mdl = mdl.Where("show_type", type_)

	// 如果选择分类，-1表示全子分类
	if subType >= 0 {
		mdl = mdl.Where("sub_show_type", subType)
	}

	// 如果选择平台，0表示全平台
	if pid > 0 {
		mdl = mdl.Where(fmt.Sprintf("JSON_CONTAINS(`platform`, '%d')", pid))
	}

	mdl = mdl.Where("status", show_model.ShowStatStandard).
		Where("is_show", show_model.ShowOn).
		Find(&res)

	fmt.Println(mdl.Error)
	ret := make([]uint64, len(res))

	for k, row := range res {
		ret[k] = row.ID
	}

	return ret
}

func showDetail(sids []uint64) []show_model.Table {
	var showRes []show_model.Table
	mdl := show_model.Model()
	mdl.Select("id", "name", "platform", "release_at").
		Where("id IN ?", sids).
		Find(&showRes)

	return showRes
}

func guduoRank(dayAt int, periodType int, st, sst int, pid int) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, 1)

	if st == int(show_model.ShowTypeMovie) && sst == int(show_model.ShowSubTypeMovieCinema) && periodType <= ListTypeMonth {
		return cinemaRank(dayAt, st, sst, pid, periodType + 1)
	}

	switch periodType {
	case ListTypeDay:
		if uint(dayAt) == time.Today() {
			res = guduoHotRankToday(st, sst, pid)
			return res
		}
		fallthrough
	case ListTypeWeek:
		fallthrough
	case ListTypeMonth:
		res = guduoHotRank(dayAt, periodType, st, sst, pid)
	case ListTypeYear:
		res = guduoPlayCountRankYear(dayAt, st, sst)
	case ListTypeAll:
		res = guduoPlayCountRank(dayAt, st, sst)
	}
	return res
}

// 进入骨朵热度排行榜
func guduoHotRankToday(st, sst, pid int) []map[string]interface{} {
	sids := showIds(st, sst, pid)
	var res []*guduo_hot_daily_model.Table
	guduo_hot_daily_model.Model().Select("show_id", "IF(custom_num != 0, custom_num, num) as num").
		Where("show_id IN ?", sids).
		Where("day_at", time.Today()).
		Order("num desc").
		Limit(50).
		Find(&res)

	resSids := make([]uint64, len(res))
	for k, row := range res {
		resSids[k] = row.ShowId
	}

	showRes := showDetail(resSids)
	ret := make([]map[string]interface{}, len(showRes))

	for k, v := range res {
		for _, row := range showRes {
			if row.ID == v.ShowId {
				ret[k] = map[string]interface{}{
					"show_id":  row.ID,
					"name":     row.Name,
					"platform": show_model.GetPlatform(row.Platform),
					"hot":      v.Num,
					"trend":    0,
					"released": (int(time.Today()) - int(row.ReleaseAt)) / (24 * 60 * 60),
				}
			}
		}
	}

	return ret
}

// 骨朵热度榜 日/周/月
func guduoHotRank(dayAt int, periodType, st, sst, pid int) []map[string]interface{} {
	var res []*guduo_hot_rank_model.Table
	mdl := guduo_hot_rank_model.Model()
	rankType := int8(0)
	switch periodType {
	case ListTypeDay:
		rankType = model_clean.CycleDaily
	case ListTypeWeek:
		rankType = model_clean.CycleWeekly
	case ListTypeMonth:
		rankType = model_clean.CycleMonthly
	}

	mdl.Select("show_id", "IF(custom_num != 0, custom_num, num) as num", "rank", "rise").
		Where("day_at", dayAt).
		Where("rank_type", rankType).
		Where("show_type", st).
		Where("sub_show_type", sst).
		Where("platform_id", pid).
		Order("num desc").
		Limit(50).
		Find(&res)

	resSids := make([]uint64, len(res))
	for k, row := range res {
		resSids[k] = row.ShowId
	}

	showRes := showDetail(resSids)
	ret := make([]map[string]interface{}, len(res))

	for k, v := range res {
		for _, row := range showRes {
			if row.ID == v.ShowId {
				ret[k] = map[string]interface{}{
					"show_id":  row.ID,
					"name":     row.Name,
					"platform": show_model.GetPlatform(row.Platform),
					"hot":      v.Num,
					"trend":    v.Rise,
					"released": (dayAt - int(row.ReleaseAt)) / (24 * 60 * 60),
				}
			}
		}
	}

	return ret
}

// 年榜
func guduoPlayCountRankYear(dayAt int, st, sst int) []map[string]interface{} {
	var res []*play_count_yearly_model.Table
	mdl := play_count_yearly_model.Model()
	mdl.Select("show_id", "num", "rise").
		Where("day_at", dayAt).
		Where("show_type", st).
		Where("sub_show_type", sst).
		Order("num desc").
		Limit(50).
		Find(&res)
	resSids := make([]uint64, len(res))
	for k, row := range res {
		resSids[k] = row.ShowId
	}

	showRes := showDetail(resSids)
	ret := make([]map[string]interface{}, len(showRes))

	for k, v := range res {
		for _, row := range showRes {
			if row.ID == v.ShowId {
				ret[k] = map[string]interface{}{
					"show_id":    row.ID,
					"name":       row.Name,
					"platform":   show_model.GetPlatform(row.Platform),
					"play_count": v.Num,
					"trend":      v.Rise,
					"released":   (dayAt - int(row.ReleaseAt)) / (24 * 60 * 60),
				}
			}
		}
	}

	return ret
}

// 总榜
func guduoPlayCountRank(dayAt, st, sst int) []map[string]interface{} {
	var res []*play_count_total_daily_model.Table
	mdl := play_count_total_daily_model.Model()
	mdl.Select("show_id", "num").
		Where("day_at", dayAt).
		Where("show_type", st).
		Where("sub_show_type", sst).
		Order("num desc").
		Limit(50).
		Find(&res)

	resSids := make([]uint64, len(res))
	for k, row := range res {
		resSids[k] = row.ShowId
	}

	showRes := showDetail(resSids)
	ret := make([]map[string]interface{}, len(showRes))
	for k, v := range res {
		for _, row := range showRes {
			if row.ID == v.ShowId {
				ret[k] = map[string]interface{}{
					"show_id":    row.ID,
					"name":       row.Name,
					"platform":   show_model.GetPlatform(row.Platform),
					"play_count": v.Num,
					"trend":      v.Rise, // concurMap[showRow.ID] - int64(k+1),
					"released":   (dayAt - int(row.ReleaseAt)) / (24 * 60 * 60),
				}
			}
		}
	}

	return ret
}

func cinemaRank(dayAt, st, sst, pid, cycle int) []map[string]interface{} {
	var res []*play_count_rank_model.Table
	play_count_rank_model.Model().Select("show_id", "num", "rank", "rise").
		Where("day_at", dayAt).
		Where("show_type", st).
		Where("sub_show_type", sst).
		Where("rank_type", cycle).
		Where("platform_id", pid).
		Order("num desc").
		Limit(50).
		Find(&res)

	resSids := make([]uint64, len(res))
	for k, row := range res {
		resSids[k] = row.ShowId
	}

	showRes := showDetail(resSids)
	ret := make([]map[string]interface{}, len(res))

	for k, v := range res {
		for _, row := range showRes {
			if row.ID == v.ShowId {
				ret[k] = map[string]interface{}{
					"show_id":  row.ID,
					"name":     row.Name,
					"platform": show_model.GetPlatform(row.Platform),
					"play_count":      v.Num,
					"trend":    v.Rise,
					"released": (dayAt - int(row.ReleaseAt)) / (24 * 60 * 60),
				}
			}
		}
	}

	return ret
}