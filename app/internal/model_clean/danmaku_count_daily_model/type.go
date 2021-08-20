package danmaku_count_daily_model


type CountRow struct {
	Num int64 `json:"num"`
}

type CountTrend struct {
	Num int64 `json:"num"`
	DayAt uint `json:"day_at"`
}
