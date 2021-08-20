package bi_shows_model

import (
	"encoding/json"
	"gorm.io/gorm"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_scrawler/show_model"
	model2 "guduo/app/migration/model"
	"guduo/app/migration/model_other"
	"guduo/pkg/db"
	"guduo/pkg/model"
	"strconv"
	"strings"
	"time"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetLoliPopMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync() {
	var tmpRes *show_model.Table
	show_model.Model().Debug().Select("id").Order("id desc").Limit(1).Find(&tmpRes)
	maxId := tmpRes.ID

	show_model.Model().Where("id > ?", maxId).Delete(nil)
	tn := show_model.Table{}.TableName()
	db.GetCrawlerMysqlConn().Exec("alter table `" + tn + "` AUTO_INCREMENT=" + strconv.FormatUint(maxId, 10) + ";")

	var res []Table
	Model().Debug().Select("*").Where("id > ?", maxId).Find(&res)


	arr := make([]*show_model.Table, 0, 500)
	arr2 := make([]uint64, 0, 500)
	for _, v := range res {
		status := int64(1)
		isShow := int8(1)
		if !v.Locked {
			isShow = 0
			status = -1
		}

		pltarr := model_other.GetPlat(v.ID)
		plts, e := json.Marshal(pltarr)
		if e != nil {
			plts = []byte("[]")
		}

		s, ss := getShowType(v.Category)
		sws := model_other.GetScriptwriter(v.ID)
		ps := model_other.GetProducer(v.ID)
		pdrs := model_other.GetPublisher(v.ID)
		staff := map[string][]string{
			"screen_writer":     sws,
			"producer":          ps,
			"producer_company":  model_other.GetCompany1(v.ID),
			"publisher":         pdrs,
			"publisher_company": model_other.GetCompany2(v.ID),
		}
		staffJson, _ := json.Marshal(staff)

		dirarr := model_other.GetDirector(v.ID)
		dirJson, e := json.Marshal(dirarr)
		if e != nil {
			dirJson = []byte("[]")
		}

		tagarr := model_other.GetShowThemeNew(v.ID)
		if len(tagarr) == 0 {
			tagarr = model_other.GetShowTheme(v.ID)
		}
		tagarr2 := make([]string, 0, 5)
		for _, vv := range tagarr {
			sss := strings.Trim(vv, " ")
			if sss != ""{
				tagarr2 = append(tagarr2, sss)
			}

		}
		tag, e := json.Marshal(tagarr2)
		if e != nil {
			tag = []byte("[]")
		}

		row := &show_model.Table{
			Name:              v.Name,
			Poster:            v.CoverImgUrl,
			ShowType:          s,
			SubShowType:       ss,
			Platform:          string(plts),
			Status:            status,
			Tag:               string(tag),
			Introduction:      v.Intro,
			Staff:             string(staffJson),
			Director:          string(dirJson),
			Length:            v.Duration,
			ReleaseAt:         model2.StrToTime(v.ReleaseDate),
			EndAt:             model2.StrToTime(v.OfflineDate),
			TotalEpisode:      v.Episode,
			IsCrawlerBaseInfo: 1,
			IsCrawlerIntro:    1,
			IsCrawlerLen:      1,
			IsCrawlerRelease:  1,
			IsShow:            isShow,
			IsSelf:            v.MadeInSelf,
			IsAdapt:           v.AdaptedStatus,
			AdaptFrom:         v.AdaptedWorksName,
			ShowStatus:        v.ReleaseStatus,
		}
		row.ID = v.ID

		arr = append(arr, row)
		arr2 = append(arr2, v.ID)
		if len(arr) >= 400 {
			show_model.Model().Unscoped().Where("id IN ?", arr2).Delete(nil)
			time.Sleep(time.Second)
			r := show_model.Model().Create(&arr)

			if r.Error != nil {
				panic("#####")
			}
			arr = arr[:0]
			arr2 = arr2[:0]
		}
	}

	if len(arr) > 0 {
		show_model.Model().Unscoped().Where("id IN ?", arr2).Delete(nil)
		show_model.Model().Create(&arr)
	}
}

func getShowType(category model.Varchar) (int64, int64) {
	switch category {
	case "AMERICAN_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesAmerica
	case "ANIME":
		return show_model.ShowTypeAmine, show_model.ShowSubTypeAmineChina
	case "ANIME(ignore)":
		return show_model.ShowTypeAmine, show_model.ShowSubTypeUnknown
	case "DOCUMENTARY":
		return show_model.ShowTypeDocumentary, show_model.ShowTypeSubDocumentary
	case "FOREIGN_KID_ANIME":
		return show_model.ShowTypeAmine, show_model.ShowSubTypeAmineForeignKid
	case "FOREIGN_KID_ANIME_MOVIE":
		return show_model.ShowTypeUnknown, show_model.ShowSubTypeAmineMovieForeignKid
	case "JAPAN_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesJP
	case "KID_ANIME_MOVIE":
		return show_model.ShowTypeAmine, show_model.ShowSubTypeAmineMovieKid
	case "KOREN_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesKR
	case "MOVIE":
		return show_model.ShowTypeMovie, show_model.ShowSubTypeMovieCinema
	case "NETWORK_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesNet
	case "NETWORK_MOVIE":
		return show_model.ShowTypeMovie, show_model.ShowSubTypeMovieNet
	case "NETWORK_VARIETY":
		return show_model.ShowTypeVariety, show_model.ShowSubTypeVarietyNet
	case "NO_DETERMINE":
		return show_model.ShowTypeUnknown, show_model.ShowSubTypeUnknown
	case "NO_DETERMINE_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeUnknown
	case "NO_DETERMINE_VARIETY":
		return show_model.ShowTypeVariety, show_model.ShowSubTypeUnknown
	case "TV_DRAMA":
		return show_model.ShowTypeSeries, show_model.ShowSubTypeSeriesTV
	case "TV_VARIETY":
		return show_model.ShowTypeVariety, show_model.ShowSubTypeVarietyTV
	}

	return show_model.ShowTypeUnknown, show_model.ShowSubTypeUnknown
}

func getPlatformId(n string) uint64 {
	if strings.Contains(n, "腾讯") {
		return constant.PlatformIdTencent
	}
	if strings.Contains(n, "优酷") {
		return constant.PlatformIdYouku
	}
	if strings.Contains(n, "芒果") {
		return constant.PlatformIdMango
	}
	if strings.Contains(n, "爱奇艺") {
		return constant.PlatformIdIqiyi
	}
	if strings.Contains(n, "bili") || strings.Contains(n, "哔哩") || strings.Contains(n, "B站") {
		return constant.PlatformIdBilibili
	}

	return 99
}
