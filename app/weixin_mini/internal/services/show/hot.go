package show

import (
	"fmt"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/article_count_daily_model"
	"guduo/app/internal/model_clean/attention_daily_model"
	"guduo/app/internal/model_clean/comment_count_daily_model"
	"guduo/app/internal/model_clean/danmaku_count_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_avg_rank_model"
	"guduo/app/internal/model_clean/guduo_hot_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_rank_model"
	"guduo/app/internal/model_clean/indicator_daily_model"
	"guduo/app/internal/model_clean/news_daily_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/time"
	"guduo/pkg/util"
)

type Hot struct {
	PlatformId int64
	PlatformStr string
	Num int64
}

// 实时骨朵热度
func CurGuduoHot(sid uint64) float64 {
	CurHot := make(map[string]interface{})
	guduoHotModel := guduo_hot_daily_model.Model()
	guduoHotModel.Debug().Select("IF(custom_num != 0, custom_num, num) as num").Where("show_id", sid).
		//Where("day_at >= ?", time.Today() - 86400 * 10). // 最多显示10天前的热度
		//Order("day_at desc").
		Where("day_at", time.Today()).
		Limit(1).
		Find(&CurHot)

	if n, ok := CurHot["num"].(float64); ok && CurHot["num"] != nil {
		return n
	}

	return .0
}

func CurGuduoRank(num float64) int64 {
	var rank int64
	guduoHotModel := guduo_hot_daily_model.Model()
	guduoHotModel.Where("IF(custom_num != 0, custom_num, num) >= ?", num).
		Where("day_at", time.Today()).
		Count(&rank)

	return rank
}

// 平均热度数据
func AvgGuduoHot(sid uint64) float64 {
	AvgHot := make(map[string]interface{})
	guduoHotModel := guduo_hot_avg_rank_model.Model()
	guduoHotModel.Select("num").
		Where("show_id", sid).
		Limit(1).
		Find(&AvgHot)
	if n, ok := AvgHot["num"].(float64); ok && AvgHot["num"] != nil {
		return util.ToFixedFloat(n, 2)
	}

	return 0.0
}

func AvgGuduoRank(sid uint64) int64 {
	//sidData := show_model.GetActiveShowsWithType()
	//sids := make([]uint64, len(sidData))
	//for k, row := range sidData {
	//	sids[k] = row.ID
	//}
	var rank *guduo_hot_avg_rank_model.Table
	guduoHotModel := guduo_hot_avg_rank_model.Model()
	guduoHotModel.Select("rank").
		Where("show_id", sid).
		Find(&rank)

	return rank.Rank
}

// 峰值热度数据
func MaxGuduoHot(sid uint64) float64 {
	MaxHot := make(map[string]interface{})
	guduoHotModel := guduo_hot_rank_model.Model()
	guduoHotModel.Select("max(IF(custom_num != 0, custom_num, num)) as num").Where("show_id", sid).
		Group("show_id").
		Find(&MaxHot)

	if n, ok := MaxHot["num"].(float64); ok && MaxHot["num"] != nil {
		return util.ToFixedFloat(n, 2)
	}

	return 0.0
}

// 峰值热度数据
func MaxGuduoHotReleaseAt(sid uint64) string {
	releaseAt := make(map[string]interface{})
	mdl := guduo_hot_rank_model.Model()
	mdl.Select("IF(custom_num != 0, custom_num, num) as num", "day_at").
		Where("show_id", sid).
		Order("num desc").
		Limit(1).
		Find(&releaseAt)
	fmt.Println(releaseAt)
	if n, ok := releaseAt["day_at"].(uint); ok && releaseAt["num"] != nil {
		return time.TimeToStr(time.LayoutYmd, n)
	}

	return ""
}


func GuduoHotTrend(sid uint64) []guduo_hot_daily_model.GuduoHotTrend {
	var trend []guduo_hot_daily_model.GuduoHotTrend
	mdl := guduo_hot_rank_model.Model()
	mdl.Select("IF(custom_num != 0, custom_num, num) as num", "day_at").
		Where("show_id", sid).
		Where("rank_type", model_clean.CycleDaily).
		Order("day_at ASC").
		Find(&trend)

	return trend
}

func SepDay(sid uint64) map[string]uint {
	var res show_model.Table

	show_model.Model().Select("release_at", "end_at" , "0 as pre_at").
		Where("id", sid).
		Find(&res)

	ret := map[string]uint {
		"release_at": res.ReleaseAt,
		"end_at": res.EndAt,
		"pre_at": 0,
	}

	return ret
}

func PlatformRank(sid uint64) map[string]int64 {
	wechatRank := getArticleRank(sid, constant.PlatformIdWeixin)
	newsRank := getNewsRank(sid, constant.PlatformIdWeixin)
	dmkRank := getDanmakuRank(sid, 0)
	baiduIndiRank := getIndicatorRank(sid, constant.PlatformIdBaidu)
	weiboIndiRank := getIndicatorRank(sid, constant.PlatformIdWeibo)
	cmtRank := getCommentRank(sid, 0)
	tiebaRank := getTiebaRank(sid)
	weiboAbout := getWeiboAbout(sid)

	hotSipderGraph := map[string]int64 {
		"wechat": wechatRank,
		"news": newsRank,
		"danmaku": dmkRank,
		"baidu": baiduIndiRank,
		"weibo": weiboIndiRank,
		"comment": cmtRank,
		"tieba": tiebaRank,
		"weibo_about": weiboAbout,
	}

	return hotSipderGraph
}

// todo 这个没有需要再爬
func getWeiboAbout(sid uint64) int64 {
	return 0
}

func getArticleRank(sid uint64, pid uint64) int64 {
	today := time.Today()
	// 微博文章排名
	num := make(map[string]interface{})
	hotModel := article_count_daily_model.Model()
	hotModel.Select("IF(custom_num != 0, custom_num, num) as num").Where("show_id", sid).
		Where("day_at", today).
		Where("platform_id", pid).
		Limit(1).
		Find(&num)

	curNum := int64(0)
	if num["num"] != nil {
		curNum, _ = num["num"].(int64)
	}
	rank := make(map[string]interface{})
	hotModel = article_count_daily_model.Model()
	hotModel.Select("count(*) as rank").
		Where("day_at", today).
		Where("platform_id", pid).
		Where("IF(custom_num != 0, custom_num, num) >= ?", curNum).
		Limit(1).
		Find(&rank)

	if r, ok := rank["rank"].(int64); ok && rank["rank"] != nil {
		return r
	}

	return 0
}


func getNewsRank(sid uint64, pid uint64) int64 {
	today := time.Today()
	// 微博文章排名
	num := make(map[string]interface{})
	hotModel := news_daily_model.Model()
	hotModel.Select("IF(custom_num != 0, custom_num, num) as num").Where("show_id", sid).
		Where("day_at", today).
		Where("platform_id", pid).
		Limit(1).
		Find(&num)

	curNum := int64(0)
	if num["num"] != nil {
		curNum, _ = num["num"].(int64)
	}
	rank := make(map[string]interface{})
	hotModel = news_daily_model.Model()
	hotModel.Select("count(*) as rank").
		Where("day_at", today).
		Where("platform_id", pid).
		Where("IF(custom_num != 0, custom_num, num) >= ?", curNum).
		Limit(1).
		Find(&rank)

	if r, ok := rank["rank"].(int64); ok && rank["rank"] != nil {
		return r
	}

	return 0
}

func getDanmakuRank(sid uint64, pid uint64) int64 {
	today := time.Today()
	// 微博文章排名
	num := make(map[string]interface{})
	hotModel := danmaku_count_daily_model.Model()
	hotModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").Where("show_id", sid).
		Where("day_at", today).
		Limit(1).
		Find(&num)

	curNum := int64(0)
	if num["num"] != nil {
		curNum, _ = num["num"].(int64)
	}

	hotModel = news_daily_model.Model()
	r := hotModel.Select("id").
		Where("day_at", today).
		Group("show_id").
		Having("sum(IF(custom_num != 0, custom_num, num)) >= ?", curNum).
		Find(nil)

	return r.RowsAffected
}


func getIndicatorRank(sid uint64, pid uint64) int64 {
	today := time.Today()
	// 微博文章排名
	num := make(map[string]interface{})
	hotModel := indicator_daily_model.Model()
	hotModel.Select("IF(custom_num != 0, custom_num, num) as num").Where("show_id", sid).
		Where("day_at", today).
		Where("platform_id", pid).
		Limit(1).
		Find(&num)

	curNum := int64(0)
	if num["num"] != nil {
		curNum, _ = num["num"].(int64)
	}
	rank := make(map[string]interface{})
	hotModel = indicator_daily_model.Model()
	hotModel.Select("count(*) as rank").
		Where("day_at", today).
		Where("platform_id", pid).
		Where("IF(custom_num != 0, custom_num, num) >= ?", curNum).
		Limit(1).
		Find(&rank)

	if r, ok := rank["rank"].(int64); ok && rank["rank"] != nil {
		return r
	}

	return 0
}


func getCommentRank(sid uint64, pid uint64) int64 {
	today := time.Today()
	// 微博文章排名
	num := make(map[string]interface{})
	hotModel := comment_count_daily_model.Model()
	hotModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").Where("show_id", sid).
		Where("day_at", today).
		Limit(1).
		Find(&num)

	curNum := int64(0)
	if num["num"] != nil {
		curNum, _ = num["num"].(int64)
	}

	hotModel = comment_count_daily_model.Model()
	r := hotModel.Select("id").
		Where("day_at", today).
		Group("show_id").
		Having("sum(IF(custom_num != 0, custom_num, num)) >= ?", curNum).
		Find(nil)

	return r.RowsAffected
}

func getTiebaRank(sid uint64) int64 {
	today := time.Today()
	// 微博文章排名
	num := make(map[string]interface{})
	hotModel := attention_daily_model.Model()
	hotModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").Where("show_id", sid).
		Where("platform_id", constant.PlatformIdBaidu).
		Where("day_at", today).
		Limit(1).
		Find(&num)

	curNum := int64(0)
	if num["num"] != nil {
		curNum, _ = num["num"].(int64)
	}

	hotModel = attention_daily_model.Model()
	r := hotModel.Select("id").
		Where("platform_id", constant.PlatformIdBaidu).
		Where("day_at", today).
		Group("show_id").
		Having("sum(IF(custom_num != 0, custom_num, num)) >= ?", curNum).
		Find(nil)

	return r.RowsAffected
}