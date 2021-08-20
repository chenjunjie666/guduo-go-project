package show_model

import (
	"encoding/json"
)

var sTypeMap = map[int64]string{
	ShowTypeSeries:  "剧集",
	ShowTypeVariety: "综艺",
	ShowTypeMovie:   "电影",
	ShowTypeAmine:   "动漫",
}

// sub show type map
var ssTypeMap = map[int64]string{
	ShowSubTypeSeriesNet:     "网剧",
	ShowSubTypeSeriesTV:      "电视剧",
	ShowSubTypeSeriesAmerica: "美剧",
	ShowSubTypeSeriesJP:      "日剧",
	ShowSubTypeSeriesKR:      "韩剧",

	// 10 开始是综艺分类
	ShowSubTypeVarietyNet: "网综",
	ShowSubTypeVarietyTV:  "电综",

	// 20 开始是电影分类
	ShowSubTypeMovieNet:    "网络大电影",
	ShowSubTypeMovieCinema: "院线电影",

	// 40 开始是动漫分类
	ShowSubTypeAmineChina: "国漫",
}

// sub show type map
var showStatusMap = map[int]string{
	ShowPlayingStatPlaying: "在播",
	ShowPlayingStatWaiting: "待播",
	ShowPlayingStatOff:     "下架",
	ShowPlayingStatInvalid: "无效",
}

var statusMap = map[int64]string{
	ShowStatReject:   "审核不通过",
	ShowStatPending:  "待审核",
	ShowStatStandard: "正常状态",
}

func GetShowTypeStr(sType int64) string {
	if name, ok := sTypeMap[sType]; ok {
		return name
	}
	return ""
}

func GetSubShowTypeStr(ssType int64) string {
	if name, ok := ssTypeMap[ssType]; ok {
		return name
	}
	return ""
}

func GetShowStatusStr(st int8) string {
	if name, ok := showStatusMap[int(st)]; ok {
		return name
	}
	return ""
}

func GetStatusStr(st int64) string {
	if name, ok := statusMap[st]; ok {
		return name
	}
	return ""
}

func GetPlatform(plt string) Platform {
	var platform Platform
	e := json.Unmarshal([]byte(plt), &platform)
	if e != nil {
		platform = make(Platform, 0)
	}

	return platform
}

func GetDirector(dirStr string) Director {
	var dir Director
	e := json.Unmarshal([]byte(dirStr), &dir)
	if e != nil {
		dir = make(Director, 0)
	}

	return dir
}

func GetTag(t string) Tag {
	var tag Tag
	e := json.Unmarshal([]byte(t), &tag)
	if e != nil {
		tag = make(Tag, 0)
	}

	return tag
}

func GetStaff(staff string, director string) Staff {
	var staffFull Staff
	var dict Director
	e := json.Unmarshal([]byte(director), &dict)
	if e != nil {
		dict = make(Director, 0)
	}

	e2 := json.Unmarshal([]byte(staff), &staffFull)
	if e2 != nil {
		staffFull = EmptyStaff()
	}
	staffFull.Director = dict
	return staffFull
}

func EmptyStaff() Staff {
	return Staff{
		make(Director, 0),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
	}
}

func GetStaffWithoutDirector(staff string) StaffWithoutDirector {
	var staffFull StaffWithoutDirector
	e2 := json.Unmarshal([]byte(staff), &staffFull)
	if e2 != nil {
		staffFull = EmptyStaffWithoutDirector()
	}
	return staffFull
}

func EmptyStaffWithoutDirector() StaffWithoutDirector {
	return StaffWithoutDirector{
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
		make([]string, 0),
	}
}
