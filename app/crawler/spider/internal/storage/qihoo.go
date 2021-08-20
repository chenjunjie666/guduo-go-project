package storage

import (
	"fmt"
	"guduo/app/internal/constant"
	actor_model "guduo/app/internal/model_scrawler/actor_model"
	indicator_actor_model "guduo/app/internal/model_scrawler/indicator_actor_model"
	indicator_model "guduo/app/internal/model_scrawler/indicator_model"
	show_model "guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"guduo/pkg/util"
)

var Qihoo = &qihoo{
	PlatformId: constant.PlatformIdQihu,
	Host:       "https://www.360.cn",
}

type qihoo struct {
	PlatformId uint64
	Host       string
}

type QihooIndicatorUrls struct {
	ActorId model.PrimaryKey
	ShowId  model.PrimaryKey
	Url     model.Varchar
}

// 获取360指数页面链接
func (q *qihoo) GetIndicatorUrl() []*QihooIndicatorUrls {
	shows := show_model.GetActiveShowsName()
	ret := make([]*QihooIndicatorUrls, len(shows))

	for k, show := range shows {
		name := util.UrlEncode(show.Name)
		url := fmt.Sprintf("https://trends.so.com/result?query=%s&period=30", name)
		ret[k] = &QihooIndicatorUrls{
			ShowId: show.Id,
			Url:    url,
		}
	}

	return ret
}

// 获取360指数页面链接
func (q *qihoo) GetIndicatorActorUrl() []*QihooIndicatorUrls {
	actors := actor_model.GetActor()
	ret := make([]*QihooIndicatorUrls, len(actors))

	for k, actor := range actors {
		name := util.UrlEncode(actor.Name)
		url := fmt.Sprintf("https://trends.so.com/result?query=%s&period=30", name)
		ret[k] = &QihooIndicatorUrls{
			ActorId: actor.Id,
			Url:     url,
		}
	}

	return ret
}

// 存储360指数信息
func (q *qihoo) StoreIndicator(i int64, jobAt uint, showId uint64) {
	indicator_model.SaveIndicator(float64(i), jobAt, showId, q.PlatformId)
}

// 存储360指数信息
func (q *qihoo) StoreIndicatorActor(i int64, jobAt uint, ActorId uint64) {
	indicator_actor_model.SaveIndicator(i, jobAt, ActorId, q.PlatformId)
}