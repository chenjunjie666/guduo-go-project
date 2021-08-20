package storage

import (
	"fmt"
	"guduo/app/internal/constant"
	actor_model "guduo/app/internal/model_scrawler/actor_model"
	attention_actor_model "guduo/app/internal/model_scrawler/attention_actor_model"
	attention_model "guduo/app/internal/model_scrawler/attention_model"
	indicator_actor_model "guduo/app/internal/model_scrawler/indicator_actor_model"
	indicator_age_model "guduo/app/internal/model_scrawler/indicator_age_model"
	indicator_gender_model "guduo/app/internal/model_scrawler/indicator_gender_model"
	indicator_model "guduo/app/internal/model_scrawler/indicator_model"
	news_model "guduo/app/internal/model_scrawler/news_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
	show_model "guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"guduo/pkg/util"
)

var Baidu = &baidu{
	PlatformId: constant.PlatformIdBaidu,
	Host:       "https://www.baidu.com/",
}

type baidu struct {
	PlatformId uint64
	Host       string
}

type GetDetailUrl struct {
	ShowId model.PrimaryKey
	Url    model.Text
}

// 获取百度贴吧页
func (b *baidu) GetDetailUrl() show_detail_model.DetailUrls {
	showIds := show_model.GetActiveShowsName()
	urls := make(show_detail_model.DetailUrls, len(showIds))
	for k, row := range showIds {
		name := util.UrlEncode(row.Name)
		url := fmt.Sprintf("https://tieba.baidu.com/f?ie=utf-8&kw=%s&fr=search", name)
		urls[k] = &struct {
			ShowId model.PrimaryKey
			Url    model.Text
		}{
			ShowId: row.Id,
			Url:    url,
		}
	}
	return urls
}

type ActorDetailUrl struct {
	ActorId model.PrimaryKey
	Url     model.Varchar
}

// 获取百度贴吧页
func (b *baidu) GetActorDetailUrl() []*ActorDetailUrl {
	actors := actor_model.GetActor()

	urls := make([]*ActorDetailUrl, len(actors))
	for k, actor := range actors {
		name := util.UrlEncode(actor.Name)
		url := fmt.Sprintf("https://tieba.baidu.com/f?ie=utf-8&kw=%s&fr=search", name)
		urls[k] = &ActorDetailUrl{
			ActorId: actor.Id,
			Url: url,
		}
	}
	return urls
}

func (b *baidu) GetNewsUrl() show_detail_model.DetailUrls {
	shows := show_model.GetActiveShowsName()
	ret := make(show_detail_model.DetailUrls, len(shows))

	for k, show := range shows {
		name := util.UrlEncode(show.Name)
		newsUrl := fmt.Sprintf("https://www.baidu.com/s?ie=utf-8&rtt=1&bsst=1&cl=2&tn=news&rsv_dl=ns_pc&word=%s", name)
		ret[k] = &struct {
			ShowId model.PrimaryKey
			Url    model.Text
		}{
			ShowId: show.Id,
			Url:    newsUrl,
		}
	}

	return ret
}

func (b *baidu) GetBaikeUrl() show_detail_model.DetailUrls {
	shows := show_model.GetActiveShowsName()

	ret := make(show_detail_model.DetailUrls, len(shows))

	for k, show := range shows {
		name := util.UrlEncode(show.Name)
		baikeUrl := fmt.Sprintf("https://baike.baidu.com/item/%s", name)
		ret[k] = &struct {
			ShowId model.PrimaryKey
			Url    model.Text
		}{
			ShowId: show.Id,
			Url:    baikeUrl,
		}
	}

	return ret
}

type BaiduIndicatorUrls struct {
	ActorId            model.PrimaryKey
	ShowId             model.PrimaryKey
	IndicatorUrl       model.Varchar
	GenAgeIndicatorUrl model.Varchar
}

func (b *baidu) GetIndicatorUrl() []*BaiduIndicatorUrls {
	shows := show_model.GetActiveShowsName()

	ret := make([]*BaiduIndicatorUrls, len(shows))

	for k, show := range shows {
		name := util.UrlEncode(show.Name)
		indiUrl := fmt.Sprintf("https://index.baidu.com/v2/main/index.html#/trend/%s?words=%s", name, name)
		genAgeUrl := fmt.Sprintf("https://index.baidu.com/v2/main/index.html#/crowd/%s?words=%s", name, name)
		ret[k] = &BaiduIndicatorUrls{
			ShowId:             show.Id,
			IndicatorUrl:       indiUrl,
			GenAgeIndicatorUrl: genAgeUrl,
		}
	}

	return ret
}

func (b *baidu) GetIndicatorActorUrl() []*BaiduIndicatorUrls {
	actors := actor_model.GetActor()

	ret := make([]*BaiduIndicatorUrls, len(actors))

	for k, actor := range actors {
		name := util.UrlEncode(actor.Name)
		indiUrl := fmt.Sprintf("https://index.baidu.com/v2/main/index.html#/trend/%s?words=%s", name, name)
		ret[k] = &BaiduIndicatorUrls{
			ActorId:      actor.Id,
			IndicatorUrl: indiUrl,
		}
	}

	return ret
}

// 存储上线时间
func (b *baidu) StoreReleaseTimeStamp(rts uint, sid uint64) {
	show_model.StoreReleaseTime(rts, sid)
	//fmt.Println(rts)
	//db.GetCrawlerMysqlConn().Create()
}

// 存储帖吧关注度和帖子数
func (b *baidu) StoreAttention(a, p int64, jobAt uint, showId uint64) {
	attention_model.StoreAttention(a, p, jobAt, showId, b.PlatformId)
}

// 存储帖吧关注度和帖子数
func (b *baidu) StoreAttentionActor(a, p int64, jobAt uint, actorId uint64) {
	attention_actor_model.StoreAttention(a, p, jobAt, actorId, b.PlatformId)
}

// 存储百度指数性别分布信息
func (b *baidu) StoreGenderRateMap(grm map[string]float64, jobAt uint, showId uint64) {
	indicator_gender_model.SaveIndicatorGender(grm, jobAt, showId, b.PlatformId)
}

// 存储百度指数性别分布信息
func (b *baidu) StoreAgeRateMap(arm map[string]float64, jobAt uint, showId uint64) {
	indicator_age_model.SaveIndicatorAge(arm, jobAt, showId, b.PlatformId)
}

// 存储百度指数信息
func (b *baidu) StoreIndicator(i int64, jobAt uint, showId uint64) {
	indicator_model.SaveIndicator(float64(i), jobAt, showId, b.PlatformId)
}

// 存储百度指数信息
func (b *baidu) StoreIndicatorActor(i int64, jobAt uint, ActorId uint64) {
	indicator_actor_model.SaveIndicator(i, jobAt, ActorId, b.PlatformId)
}

// 存储新闻资讯数
func (b *baidu) StoreNewsNum(nn int64, ja uint, sid uint64) {
	news_model.SaveNewsNum(nn, ja, sid, b.PlatformId)
	//fmt.Println(nn)
	//db.GetCrawlerMysqlConn().Create()
}

//func (b baidu) StorePostNum(pn int64) {
//	fmt.Println(pn)
//}
