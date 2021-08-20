package storage

import (
	"fmt"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/app/internal/model_scrawler/article_num_actor_model"
	"guduo/app/internal/model_scrawler/article_num_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"guduo/pkg/util"
)

var Wechat = &wechat{
	PlatformId: constant.PlatformIdWeixin,
	Host:       "https://weixin.sogou.com/",
}

type wechat struct {
	PlatformId uint64
	Host       string
}

type GetArticleNumUrlReturn struct {
	ID model.PrimaryKey
	Url model.Varchar
}

// 获取微信文章数详情页面链接
func (w *wechat) GetArticleNumUrl() []GetArticleNumUrlReturn {
	showsName := show_model.GetActiveShowsName()

	urls := make([]GetArticleNumUrlReturn, 0, 50)
	for _, row := range showsName {
		nameEncode := util.UrlEncode(row.Name)
		u := fmt.Sprintf("https://weixin.sogou.com/weixin?type=2&query=%s", nameEncode)
		ret := GetArticleNumUrlReturn{
			ID:   row.Id,
			Url: u,
		}
		urls = append(urls, ret)
	}

	return urls
}

// 获取微信文章数详情页面链接
func (w *wechat) GetArticleNumActorUrl() []GetArticleNumUrlReturn {
	actorName := actor_model.GetActor()

	urls := make([]GetArticleNumUrlReturn, 0, 50)
	for _, row := range actorName {
		nameEncode := util.UrlEncode(row.Name)
		u := fmt.Sprintf("https://weixin.sogou.com/weixin?type=2&query=%s", nameEncode)
		ret := GetArticleNumUrlReturn{
			ID:   row.Id,
			Url: u,
		}
		urls = append(urls, ret)
	}

	return urls
}

// 存储文章数
func (w *wechat) StoreArticleNum(an int64, showId uint64, jobAt uint) {
	article_num_model.SaveArticleNum(an, jobAt, showId, w.PlatformId)
}

// 存储文章数
func (w *wechat) StoreArticleNumActor(an int64, actorId uint64, jobAt uint) {
	article_num_actor_model.SaveArticleNum(an, jobAt, actorId, w.PlatformId)
}
