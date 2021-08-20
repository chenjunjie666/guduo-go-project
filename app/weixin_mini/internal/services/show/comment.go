package show

import (
	"guduo/app/internal/model_clean/comment_count_trend_daily_model"
)

func TotalCommentCount(sid uint64, pid ...uint64) int64 {
	cmtCount := make(map[string]interface{})

	cmtModel := comment_count_trend_daily_model.Model()

	if len(pid) != 0 {
		cmtModel.Where("platform_id IN ?", pid)
	}

	cmtModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").
		Where("show_id", sid).
		Find(&cmtCount)

	if n, ok := cmtCount["num"].(int64); cmtCount["num"] != nil && ok {
		return n
	}

	return 0
}

func DayCommentCount(sid uint64, day uint,  pid ...uint64) int64 {
	cmtCount := make(map[string]interface{})

	cmtModel := comment_count_trend_daily_model.Model()

	if len(pid) != 0 {
		cmtModel.Where("platform_id IN ?", pid)
	}

	cmtModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").
		Where("show_id", sid).
		Where("day_at", day).
		Find(&cmtCount)

	if n, ok := cmtCount["num"].(int64); cmtCount["num"] != nil && ok {
		return n
	}

	return 0
}

func DayCommentRank(sid uint64, day uint, num int64, pid ...uint64) int64 {
	cmtModel := comment_count_trend_daily_model.Model()
	if len(pid) != 0 {
		cmtModel = cmtModel.Where("platform_id IN ?", pid)
	}

	var cnt int64
	cmtModel.Select("id", "show_id").
		Where("day_at", day).
		Group("show_id").
		Having("sum(IF(custom_num != 0, custom_num, num)) >= ?", num).
		Count(&cnt)

	return cnt
}


func CommentCountTrend(sid uint64) []comment_count_trend_daily_model.CountTrend {
	var trend []comment_count_trend_daily_model.CountTrend
	cmtModel := comment_count_trend_daily_model.Model()
	cmtModel.Select("SUM(IF(custom_num != 0, custom_num, num)) as num", "day_at").
		Where("show_id", sid).
		Group("day_at").
		Order("day_at ASC").
		Find(&trend)

	return trend
}