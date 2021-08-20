package show

import (
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/article_content_current_model"
	"guduo/app/internal/model_clean/article_count_trend_daily_model"
	"guduo/app/internal/model_scrawler/show_model"
)

func DayArticleCount(sid uint64, day uint,  pid ...uint64) int64 {
	cnt := make(map[string]interface{})

	articleModel := article_count_trend_daily_model.Model()

	if len(pid) != 0 {
		articleModel = articleModel.Where("platform_id IN ?", pid)
	}

	articleModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").
		Where("show_id", sid).
		Where("day_at", day).
		Find(&cnt)

	if n, ok := cnt["num"].(int64); cnt["num"] != nil && ok {
		return n
	}

	return 0
}

func DayArticleRank(sid uint64, day uint, num int64, pid ...uint64) (int64, string) {
	tp := show_model.GetShowInfo(sid)
	shows := show_model.GetActiveShowsByType(tp.ShowType, tp.SubShowType)
	showTypeStr := show_model.GetSubShowTypeStr(tp.SubShowType)
	if len(shows) == 0 {
		return 0, ""
	}
	sids := make([]uint64, len(shows))
	for k, row := range shows {
		sids[k] = row.ID
	}
	articleModel := article_count_trend_daily_model.Model()
	if len(pid) != 0 {
		articleModel= articleModel.Where("platform_id IN ?", pid)
	}

	var cnt int64
	articleModel.Select("id").
		Where("day_at", day).
		Where("show_id IN ?", sids).
		Group("show_id").
		Having("sum(IF(custom_num !=0, custom_num, num)) >= ?", num).
		Count(&cnt)

	return cnt, showTypeStr
}


func DayArticleNumTrend(sid uint64) []article_count_trend_daily_model.CountTrend {
	trend := make([]article_count_trend_daily_model.CountTrend, 0, 100)
	cntModel := article_count_trend_daily_model.Model()
	cntModel.Select("SUM(IF(custom_num != 0, custom_num, num)) as num", "day_at").
		Where("show_id", sid).
		Where("platform_id", constant.PlatformIdWeibo).
		Group("day_at").
		Order("day_at ASC").
		Find(&trend)


	return trend
}



func CurHotArticle(sid uint64) []map[string]interface{}{
	var res []article_content_current_model.Table
	
	articleModel := article_content_current_model.Model()
	articleModel.Select("author", "content", "publish_at", "forward").
		Where("platform_id", constant.PlatformIdWeibo).
		Where("show_id", sid).
		Order("forward desc").
		Limit(5).
		Find(&res)
	
	ret := make([]map[string]interface{}, len(res))
	for k, row := range res {
		ret[k] = map[string]interface{}{
			"author": row.Author,
			"content": row.Content,
			"publish_at": row.PublishAt,
			"forward": row.Forward,
		}
	}

	return ret
}