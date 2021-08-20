package storage

import (
	"guduo/app/internal/constant"
	rating_model "guduo/app/internal/model_scrawler/rating_model"
	short_comment_model "guduo/app/internal/model_scrawler/short_comment_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
	show_model "guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/model"
)

var Douban = &douban{
	PlatformId: constant.PlatformIdDouban,
	Host:       "https://www.douban.com",
}

type douban struct {
	PlatformId uint64
	Host       string
}

func (d *douban) GetNeedFetchLengthUrl() show_detail_model.DetailUrls {
	f := show_model.GetActiveShows()
	sids := make([]model.PrimaryKey, len(f))
	for k, v := range f {
		sids[k] = v
	}

	urls := show_detail_model.GetDetailUrl(d.PlatformId, sids)
	return urls
}

// 获取豆瓣电影详情页面链接
func (d *douban) GetDetailUrl() show_detail_model.DetailUrls {
	f := show_model.GetActiveShows()
	sids := make([]model.PrimaryKey, len(f))
	for k, v := range f {
		sids[k] = v
	}

	urls := show_detail_model.GetDetailUrl(d.PlatformId, sids)
	return urls
	//urls := []string{
	//	"https://movie.douban.com/subject/27663998/",
	//}
	//return urls
}

// TODO
// 存储评分
func (d *douban) StoreRatingNum(rn float64, ja uint, sid uint64) {
	rating_model.SaveRatingNum(rn, ja, sid, d.PlatformId)
	//db.GetCrawlerMysqlConn().Create()
}

// 存储短评数
func (d *douban) StoreShortCommentCount(scc int64, ja uint, sid uint64) {
	short_comment_model.SaveShortCommentCount(scc, ja, sid, d.PlatformId)
}

// 存储片长
func (d *douban) StoreLength(len string, sid uint64) {

	storeLength(len, sid)
}
