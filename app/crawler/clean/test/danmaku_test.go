package test

import (
	"fmt"
	"guduo/app/crawler/clean/internal/core"
	danmaku_word_cloud_daily_model "guduo/app/internal/model_clean/danmaku_word_cloud_daily_model"
	danmaku_model "guduo/app/internal/model_scrawler/danmaku_model"
	"testing"
)

func TestGetDanmaku(t *testing.T) {
	core.Init()
	res := danmaku_model.GetDanmaku(1)
	fmt.Println(res)
}

func TestSaveDanmaku(t *testing.T)  {
	core.Init()
	json := "[{},{},{}]"
	picBase64 := "adsasdasd"
	danmaku_word_cloud_daily_model.SaveWordCloud(json, picBase64, 1000000000, 1)
}