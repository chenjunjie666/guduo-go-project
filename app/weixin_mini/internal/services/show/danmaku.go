package show

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/psykhi/wordclouds"
	"guduo/app/internal/model_clean/danmaku_count_daily_model"
	"guduo/app/internal/model_clean/danmaku_count_trend_daily_model"
	"guduo/app/internal/model_clean/danmaku_word_cloud_daily_model"
	"image/color"
	"image/jpeg"
	"os"
)

func TotalDanmakuCount(sid uint64, pid ...uint64) int64 {
	dmkCount := make(map[string]interface{})

	dmkModel := danmaku_count_daily_model.Model()

	if len(pid) != 0 {
		dmkModel = dmkModel.Where("platform_id IN ?", pid)
	}

	// 总弹幕就是  danmaku_count_daily 表今天的值
	dmkModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").
		Where("show_id", sid).
		Find(&dmkCount)

	if n, ok := dmkCount["num"].(int64); dmkCount["num"] != nil && ok {
		return n
	}

	return 0
}

func DayDanmakuCount(sid uint64, day uint,  pid ...uint64) int64 {
	dmkCount := make(map[string]interface{})

	dmkModel := danmaku_count_trend_daily_model.Model()

	if len(pid) != 0 {
		dmkModel.Where("platform_id IN ?", pid)
	}

	dmkModel.Select("sum(IF(custom_num != 0, custom_num, num)) as num").
		Where("show_id", sid).
		Where("day_at", day).
		Find(&dmkCount)

	if n, ok := dmkCount["num"].(int64); dmkCount["num"] != nil && ok {
		return n
	}

	return 0
}

func DayDanmakuRank(sid uint64, day uint, num int64, pid ...uint64) int64 {
	dmkModel := danmaku_count_trend_daily_model.Model()
	if len(pid) != 0 {
		dmkModel.Where("platform_id IN ?", pid)
	}

	var cnt int64
	dmkModel.Where("show_id", sid).
		Where("day_at", day).
		Group("show_id").
		Having("sum(IF(custom_num != 0, custom_num, num)) >= ?", num).
		Count(&cnt)

	return cnt
}


func DanmakuCountTrend(sid uint64) []danmaku_count_trend_daily_model.CountTrend {
	var trend []danmaku_count_trend_daily_model.CountTrend
	dmkModel := danmaku_count_trend_daily_model.Model()
	dmkModel.Select("SUM(IF(custom_num != 0, custom_num, num)) as num", "day_at").
		Where("show_id", sid).
		Group("day_at").
		Order("day_at ASC").
		Find(&trend)

	return trend
}


func WordCloud(sid uint64) string {
	var res []danmaku_word_cloud_daily_model.Table

	mdl := danmaku_word_cloud_daily_model.Model()
	mdl.Select("word", "weight").
		Where("show_id", sid).
		Order("day_at desc").
		Order("weight desc").
		Limit(50).
		Find(&res)
	fmt.Println(res)
	ws := make(map[string]int)

	weight := 100
	for _, row := range res {
		words := "\n  " + row.Word + "  \n"
		ws[words] = weight
		if weight > 40 {
			weight = int(float64(weight) * 0.92)
		}
	}


	return genPic(ws)
}


var DefaultColors = []color.RGBA{
	{0xed, 0xcc, 0x4c, 0xff},
	{0x8b, 0xc8, 0xd6, 0xff},
	{0x59, 0x61, 0xa6, 0xff},
	{0xe3, 0x8e, 0x3b, 0xff},
	{0xa5, 0xc6, 0x4f, 0xff},
	{0x67, 0xae, 0xe0, 0xff},
	{0xd0, 0x8e, 0xb6, 0xff},
	{0xdf, 0x83, 0x65, 0xff},
}

func genPic(weightMap map[string]int) string {
	colors := make([]color.Color, 0)
	for _, c := range DefaultColors {
		colors = append(colors, c)
	}

	path := "/Users/zhaokun/Documents/go/src/guduo-crawler/build/black.ttf"
	if !ttfExists(path) {
		path = "/Users/zhaokun/Documents/Golang/src/guduo/build/black.ttf"
	}

	if !ttfExists(path) {
		path = "black.ttf"
	}

	w := wordclouds.NewWordcloud(
		weightMap,
		wordclouds.FontFile(path),
		wordclouds.FontMaxSize(120),
		wordclouds.FontMinSize(60),
		wordclouds.Colors(colors),
		wordclouds.Width(1500),
		wordclouds.Height(750),
	)

	img := w.Draw()

	emptyBuff := bytes.NewBuffer(nil)
	jpeg.Encode(emptyBuff, img, nil)
	dist := base64.StdEncoding.EncodeToString(emptyBuff.Bytes())
	picB64 := dist
	return picB64
}

func ttfExists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}