package daily

import (
	"guduo/app/crawler/clean/internal/core"
	"guduo/pkg/time"
	"testing"
)

func TestGuduoDanmakuWordCloud(t *testing.T) {
	core.Init()


	Today = time.Today()

	Yesterday = Today - 86400

	beforeYesterday = Yesterday - 86400


	//guduo_hot_rank_model.Model().
	//	Where("day_at", Yesterday).
	//	Where("rank_type", model_clean.CycleDaily).
	//	Delete(nil)
	guduoWordCloud(22270)
}