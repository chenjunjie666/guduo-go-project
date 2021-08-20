package storage

import (
	log "github.com/sirupsen/logrus"
	"guduo/app/internal/constant"
	comment_count_model "guduo/app/internal/model_scrawler/comment_count_model"
	"guduo/app/internal/model_scrawler/danmaku_count_model"
	danmaku_model "guduo/app/internal/model_scrawler/danmaku_model"
	show_detail_model "guduo/app/internal/model_scrawler/show_detail_model"
)

var Bilibili = &bilibili{
	PlatformId: constant.PlatformIdBilibili,
	Host:       "https://www.bilibili.com",
}

type bilibili struct {
	PlatformId uint64
	Host       string
}

// 获取Bilibili电影详情页面链接
func (b *bilibili) GetDetailUrl() show_detail_model.DetailUrls {
	//showIds := show_model.GetActiveShows()
	//urls := show_detail_model.GetDetailUrl(b.PlatformId, showIds)
	urls := show_detail_model.GetDetailUrlNew(b.PlatformId)

	log.Info("bilibili一共", len(urls), "部剧综")
	return urls
	//urls := []string{
	//	"https://www.bilibili.com/bangumi/play/ep374527?spm_id_from=333.851.b_62696c695f7265706f72745f616e696d65.9",
	//	"https://www.bilibili.com/video/BV1Dv411873X",
	//}
	//return urls
}

// 存储评论数
func (b *bilibili) StoreCommentCount(cc int64, jobAt uint, showId uint64) {
	comment_count_model.SaveCommentCount(cc, jobAt, showId, b.PlatformId)
}

// 存储弹幕内容
func (b *bilibili) StoreDanmakuContent(cts []string, ja uint, sid uint64) {
	danmaku_model.SaveDanmaku(cts, ja, sid, b.PlatformId)
}

func (b bilibili) StoreDanmakuCount(dc int64, ja uint, sid uint64) {
	danmaku_count_model.SaveDanmakuCount(dc, ja, sid, b.PlatformId)
}
