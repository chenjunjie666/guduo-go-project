package show

import (
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/rating_daily_model"
	"guduo/pkg/time"
	"guduo/pkg/util"
)

func CurDoubanRating(sid uint64) float64 {
	CurRating := make(map[string]interface{})
	ratingModel := rating_daily_model.Model()
	ratingModel.Select("IF(custom_rating != 0, custom_rating, rating) as rating").
		Where("show_id", sid).
		Where("day_at", time.Today()).
		Where("platform_id", constant.PlatformIdDouban).
		Limit(1).
		Find(&CurRating)

	if n, ok := CurRating["rating"].(float64); ok && CurRating["rating"] != nil {
		return util.ToFixedFloat(n, 2)
	}

	return 0
}


func MaxDoubanRating(sid uint64) float64 {
	var MaxRating rating_daily_model.RatingRow
	ratingModel := rating_daily_model.Model()
	ratingModel.Select("max(IF(custom_rating != 0, custom_rating, rating)) as rating").
		Where("show_id", sid).
		Where("platform_id", constant.PlatformIdDouban).
		Limit(1).
		Find(&MaxRating)
	return util.ToFixedFloat(MaxRating.Rating, 2)
}


func DoubanRatingTrend(sid uint64) []rating_daily_model.RatingTrend {
	trend := make([]rating_daily_model.RatingTrend, 0, 100)
	ratingModel := rating_daily_model.Model()
	ratingModel.Select("IF(custom_rating != 0, custom_rating, rating) as rating", "day_at").
		Where("show_id", sid).
		Where("rating > ?", 0).
		Where("platform_id", constant.PlatformIdDouban).
		Order("day_at ASC").
		Find(&trend)


	return trend
}
