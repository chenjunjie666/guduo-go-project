package storage

import (
	log "github.com/sirupsen/logrus"
	"guduo/app/internal/constant"
	comment_count_model "guduo/app/internal/model_scrawler/comment_count_model"
	"guduo/app/internal/model_scrawler/danmaku_count_model"
	"guduo/app/internal/model_scrawler/danmaku_model"
	hot_model "guduo/app/internal/model_scrawler/hot_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
)

var Iqiyi = &iqiyi{
	PlatformId: constant.PlatformIdIqiyi,
	Host:       "https://www.iqiyi.com",
}

type iqiyi struct {
	PlatformId uint64
	Host       string
}

// 获取爱奇艺详情页url
func (i *iqiyi) GetDetailUrl() show_detail_model.DetailUrls {
	//showIds := show_model.GetActiveShows()
	//urls := show_detail_model.GetDetailUrl(i.PlatformId, showIds)
	urls := show_detail_model.GetDetailUrlNew(i.PlatformId)

	log.Info("爱奇艺共", len(urls), "个剧需要爬取")
	//urls := []string{
	//	"https://www.iqiyi.com/v_260uudpmizo.html",
	//	"https://www.iqiyi.com/v_19ruzj8gv0.html",
	//}
	return urls
}

func (i *iqiyi) GetNeedFetchBaseInfoUrl() show_detail_model.DetailUrls {
	return getNeedFetchBaseInfoUrl(i.PlatformId)
}

// 存储热度
func (i *iqiyi) StoreHot(hot int64, jobAt uint, showId uint64) {
	hot_model.SaveHotCount(hot, jobAt, showId, i.PlatformId)
}

// 存储评论数量
func (i *iqiyi) StoreCommentCount(cc int64, jobAt uint, showId uint64) {
	comment_count_model.SaveCommentCount(cc, jobAt, showId, i.PlatformId)
}

// 存储艺人信息
func (i *iqiyi) StoreBaseInfoMap(bim map[string]string, showId uint64) {
	storeBaseInfo(bim, showId)
}

// 存储弹幕数
func (i *iqiyi) StoreDanmakuCount(dc int64, ja uint, sid uint64) {
	danmaku_count_model.SaveDanmakuCount(dc, ja, sid, i.PlatformId)
}

// 存储弹幕内容
func (i *iqiyi) StoreDanmakuContent(cts []string, ja uint, sid uint64) {
	danmaku_model.SaveDanmaku(cts, ja, sid, i.PlatformId)
}
