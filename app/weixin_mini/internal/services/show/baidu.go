package show

import (
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/indicator_age_daily_model"
	indicator_gender_daily_model "guduo/app/internal/model_clean/indicator_gender_daily_model"
	"guduo/pkg/util"
)

func AgeIndicator(sid uint64) []map[string]interface{} {
	var dRes *indicator_age_daily_model.Table
	indicator_age_daily_model.Model().Select("day_at").Where("show_id", sid).
		Where("platform_id", constant.PlatformIdBaidu).
		Order("day_at desc").
		Limit(1).
		Find(&dRes)

	if dRes.DayAt == 0 {
		return make([]map[string]interface{}, 0)
	}

	dayAt := dRes.DayAt

	var res []indicator_age_daily_model.Table

	mdl := indicator_age_daily_model.Model()
	mdl.Select("age_from", "age_to", "IF(custom_rating != 0, custom_rating, rating) as rating").
		Where("show_id", sid).
		Where("platform_id", constant.PlatformIdBaidu).
		Where("day_at", dayAt).
		Order("age_from asc").
		Find(&res)

	ret := make([]map[string]interface{}, len(res))

	for k, row := range res {
		ret[k] = map[string]interface{}{
			"age_from": row.AgeFrom,
			"age_to":   row.AgeTo,
			"rating":   util.ToFixedFloat(row.Rating, 2),
		}
	}

	return ret
}

func GenderIndicator(sid uint64) map[string]float64 {
	var res indicator_gender_daily_model.Table

	mdl := indicator_gender_daily_model.Model()
	mdl.Select("male_rating", "female_rating").
		Where("show_id", sid).
		Where("platform_id", constant.PlatformIdBaidu).
		Order("day_at desc").
		Limit(1).
		Find(&res)

	ret := map[string]float64{
		"male_rating":   util.ToFixedFloat(res.MaleRating, 2),
		"female_rating": util.ToFixedFloat(res.FemaleRating, 2),
	}

	return ret
}
