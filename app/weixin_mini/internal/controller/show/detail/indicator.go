package detail

import (
	"github.com/gin-gonic/gin"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/rating_daily_model"
	"guduo/app/weixin_mini/internal/hepler/request"
	"guduo/app/weixin_mini/internal/hepler/resp"
	"guduo/app/weixin_mini/internal/services/show"
	"guduo/pkg/time"
)

// 全网热度、豆瓣口碑、弹幕评论、弹幕词云、微博、热门微博、受众分析

// 全网热度
func NetHot(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}

	// 实时热度、平均热度、热度峰值
	curHot := show.CurGuduoHot(sid)
	curRank := show.CurGuduoRank(curHot)
	avgHot := show.AvgGuduoHot(sid)
	avgRank := show.AvgGuduoRank(sid)
	MaxHot := show.MaxGuduoHot(sid)
	ReleaseAt := show.MaxGuduoHotReleaseAt(sid)

	// 排名，蛛网图
	hotTrend := show.GuduoHotTrend(sid)

	sepDay := show.SepDay(sid)
	if len(hotTrend) > 0 {
		if sepDay["release_at"] == 0 {
			sepDay["release_at"] = hotTrend[0].DayAt
		}
		if sepDay["end_at"] == 0 {
			sepDay["end_at"] = hotTrend[len(hotTrend)-1].DayAt
		}
	}

	ret := map[string]interface{}{
		"cur":      curHot,
		"cur_rank": curRank,
		"avg":      avgHot,
		"avg_rank": avgRank,
		"max":      MaxHot,
		"max_at":   ReleaseAt,
		"trend":    hotTrend,
		"sep_day":  sepDay,
	}

	resp.Success(c, ret)
}

// 豆瓣评分人数 没有爬取
// 当前豆瓣评分，历史最高评分，过去30天的评分走势
func Douban(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}

	curRating := show.CurDoubanRating(sid)
	maxRating := show.MaxDoubanRating(sid)
	RatingTrend := make([]rating_daily_model.RatingTrend, 0, 100)
	RatingTrend = show.DoubanRatingTrend(sid)

	allZero := true
	for _, v := range RatingTrend {
		if v.DayAt > 0 {
			allZero = false
			break
		}
	}

	if allZero {
		RatingTrend = make([]rating_daily_model.RatingTrend, 0)
	}

	ret := map[string]interface{}{
		"cur":   curRating,
		"max":   maxRating,
		"trend": RatingTrend,
	}

	resp.Success(c, ret)
}

// 弹幕总数，昨日弹幕数，昨日排名（子分类下）以及趋势
// 评论总数，昨日评论数，昨日排名（子分类下）以及趋势
func DanmakuComment(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}
	// 弹幕相关
	totalDmkCount := show.TotalDanmakuCount(sid)
	yesterdayDmkCount := show.DayDanmakuCount(sid, time.Today()-86400)
	yesterdayDmkRank := show.DayDanmakuRank(sid, time.Today()-86400, yesterdayDmkCount)

	beforeDayDmkCount := show.DayDanmakuCount(sid, time.Today()-86400*2)
	beforeDayDmkRank := show.DayDanmakuRank(sid, time.Today()-86400*2, beforeDayDmkCount)

	// 弹幕排名变化趋势
	dmkRankTrend := 0
	if yesterdayDmkRank > beforeDayDmkRank {
		dmkRankTrend = 1
	} else if yesterdayDmkRank < beforeDayDmkRank {
		dmkRankTrend = -1
	}
	dmkCountTrend := show.DanmakuCountTrend(sid)
	totalDanmak := int64(0)
	for _, v := range dmkCountTrend {
		totalDanmak += v.Num
	}
	avgDanmaku := int64(0)
	if int64(len(dmkCountTrend)) > 0 {
		avgDanmaku = totalDanmak / int64(len(dmkCountTrend))
	}

	// 品论相关
	totalCmtCount := show.TotalCommentCount(sid)
	yesterdayCmtCount := show.DayCommentCount(sid, time.Today()-86400)
	yesterdayCmtRank := show.DayCommentRank(sid, time.Today()-86400, yesterdayCmtCount)

	beforeDayCmtCount := show.DayCommentCount(sid, time.Today()-86400*2)
	beforeDayCmtRank := show.DayCommentRank(sid, time.Today()-86400*2, beforeDayCmtCount)
	// 弹幕排名变化趋势
	cmtRankTrend := 0
	if yesterdayCmtRank > beforeDayCmtRank {
		cmtRankTrend = 1
	} else if yesterdayCmtRank < beforeDayCmtRank {
		cmtRankTrend = -1
	}
	cmtCountTrend := show.CommentCountTrend(sid)

	ret := map[string]map[string]interface{}{
		"danmaku": {
			"total_count":          totalDmkCount,
			"yesterday_count":      yesterdayDmkCount,
			"avg_count":            avgDanmaku,
			"yesterday_rank":       yesterdayDmkRank,
			"yesterday_rank_trend": dmkRankTrend,
			"count_trend":          dmkCountTrend,
		},
		"comment": {
			"total_count":          totalCmtCount,
			"yesterday_count":      yesterdayCmtCount,
			"yesterday_rank":       yesterdayCmtRank,
			"yesterday_rank_trend": cmtRankTrend,
			"count_trend":          cmtCountTrend,
		},
	}

	resp.Success(c, ret)
}

// 弹幕词云
func WordCloud(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}

	picB64 := show.WordCloud(sid)

	ret := map[string]string{
		"img": picB64,
	}

	resp.Success(c, ret)
}

// 昨日相关微博数，昨日排名（子分类下）以及趋势
func Weibo(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}
	// 弹幕相关
	yesterdayArticleCount := show.DayArticleCount(sid, time.Today()-86400, constant.PlatformIdWeibo)
	yesterdayArticleRank, subType := show.DayArticleRank(sid, time.Today()-86400, yesterdayArticleCount, constant.PlatformIdWeibo)
	beforeDayArticleCount := show.DayArticleCount(sid, time.Today()-86400*2, constant.PlatformIdWeibo)
	beforeDayArticleRank, _ := show.DayArticleRank(sid, time.Today()-86400*2, beforeDayArticleCount, constant.PlatformIdWeibo)

	// 弹幕排名变化趋势
	articleRankTrend := 0
	if yesterdayArticleRank > beforeDayArticleRank {
		articleRankTrend = 1
	} else if yesterdayArticleRank < beforeDayArticleRank {
		articleRankTrend = -1
	}

	articleCountTrend := show.DayArticleNumTrend(sid)
	ret := map[string]map[string]interface{}{
		"weibo": {
			"yesterday_count":      yesterdayArticleCount,
			"yesterday_rank":       yesterdayArticleRank,
			"yesterday_rank_trend": articleRankTrend,
			"trend":                articleCountTrend,
			"sub_type_str":         subType,
		},
	}

	resp.Success(c, ret)
}

// 当前热门微博
func WeiboHot(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}

	ret := show.CurHotArticle(sid)

	resp.Success(c, ret)
}

// 受众分析-性别分布，年龄分布
func Analysis(c *gin.Context) {
	sid := request.GetShowId(c)
	if sid <= 0 {
		resp.Fail(c, "非法请求")
		return
	}

	age := show.AgeIndicator(sid)
	gender := show.GenderIndicator(sid)

	ret := map[string]interface{}{
		"age":    age,
		"gender": gender,
	}

	resp.Success(c, ret)
}
