package guduo_hot_daily_model

import "guduo/pkg/model"

type GuduoHotTrend struct {
	Num model.Float `json:"num"`
	DayAt model.SecondTimeStamp `json:"day_at"`
}
