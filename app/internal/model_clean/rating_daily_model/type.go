package rating_daily_model

type RatingRow struct {
	Rating float64 `json:"rating"`
}

type RatingTrend struct {
	Rating float64 `json:"rating"`
	DayAt float64 `json:"day_at"`
}
