package storage

import (
	"fmt"
	"guduo/app/internal/constant"
	actor_model "guduo/app/internal/model_scrawler/actor_model"
	article_content_model "guduo/app/internal/model_scrawler/article_content_model"
	"guduo/app/internal/model_scrawler/article_num_model"
	fans_model "guduo/app/internal/model_scrawler/fans_model"
	indicator_actor_model "guduo/app/internal/model_scrawler/indicator_actor_model"
	indicator_model "guduo/app/internal/model_scrawler/indicator_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
	show_model "guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
	"guduo/pkg/util"
)

var Weibo = &weibo{
	PlatformId: constant.PlatformIdWeibo,
	Host:       "https://www.weibo.com",
}

type weibo struct {
	PlatformId uint64
	Host       string
}

func (w *weibo) GetMainPageUrl() []string {
	return nil
}

type GetFetchArticleUrlReturn struct {
	ID  model.PrimaryKey
	Url model.Varchar
}

func (w *weibo) GetDetailUrl() show_detail_model.DetailUrls {
	sids := show_model.GetActiveShows()
	urls := show_detail_model.GetDetailUrl(w.PlatformId, sids)
	return urls
}

// 获取抓取微博文章的链接
func (w *weibo) GetFetchArticleUrl() []GetFetchArticleUrlReturn {
	f := show_model.GetActiveShowsName()

	urls := make([]GetFetchArticleUrlReturn, 0, 50)
	for _, row := range f {
		nameEncode := util.UrlEncode(row.Name)
		u := fmt.Sprintf("https://s.weibo.com/weibo/%s?topnav=1&wvr=6&b=1", nameEncode)

		ret := GetFetchArticleUrlReturn{
			ID:  row.Id,
			Url: u,
		}
		urls = append(urls, ret)
	}

	return urls
}

// 获取微博指数所需要的剧集名
func (w *weibo) GetIndicatorUrl() []*show_model.ShowName {
	shows := show_model.GetActiveShowsName()

	return shows
	//urls := []string{
	//	"https://data.weibo.com/index/newindex?visit_type=trend&wid=1060000007626",
	//}
	//return urls
}

// 获取微博指数所需要的剧集名
func (w *weibo) GetIndicatorActorUrl() []*actor_model.ActorsName {
	actors := actor_model.GetActor()

	return actors
	//urls := []string{
	//	"https://data.weibo.com/index/newindex?visit_type=trend&wid=1060000007626",
	//}
	//return urls
}

// 存储微博粉丝数
func (w *weibo) StoreFansCount(fc int64, jobAt uint, showId uint64) {
	fans_model.SaveFansCount(fc, jobAt, showId, w.PlatformId)
}

// 存储微博指数
func (w *weibo) StoreIndicator(in int64, jobAt uint, showId uint64) {
	indicator_model.SaveIndicator(float64(in), jobAt, showId, w.PlatformId)
}

// 存储微博指数
func (w *weibo) StoreIndicatorActor(in int64, jobAt uint, actorId uint64) {
	indicator_actor_model.SaveIndicator(in, jobAt, actorId, w.PlatformId)
}


func (w *weibo) StoreArticleContent(uid, name, _time, content string, forward int64, jobAt uint, showId uint64) int64 {
	return article_content_model.SaveContent(uid, name, _time, content, forward, jobAt, showId, w.PlatformId)
}

func (w weibo) StoreArticleNum(tc int64, jobAt uint, showId uint64)  {
	article_num_model.SaveArticleNum(tc, jobAt, showId, w.PlatformId)
}
