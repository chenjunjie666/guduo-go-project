// 骨朵热度
package free

import (
	"guduo/app/crawler/clean/task"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean/article_count_trend_daily_model"
	"guduo/app/internal/model_clean/attention_trend_daily_model"
	"guduo/app/internal/model_clean/comment_count_trend_daily_model"
	"guduo/app/internal/model_clean/danmaku_count_trend_daily_model"
	"guduo/app/internal/model_clean/guduo_hot_daily_model"
	"guduo/app/internal/model_clean/indicator_daily_model"
	"guduo/app/internal/model_clean/news_trend_daily_model"
	"guduo/app/internal/model_scrawler/show_model"
	"guduo/pkg/time"
	"guduo/pkg/util"
	"math"
	"math/rand"
	time2 "time"
)

func GuduoHotHandle() {
	guduoHotCalc()
}

var factors = map[int64][7]float64{
	show_model.ShowSubTypeVarietyTV:  {0.4, 0, 0.1, 0.15, 0.05, 0.3, 0},
	show_model.ShowSubTypeVarietyNet: {0.4, 0, 0.1, 0.15, 0.05, 0.3, 0},
	show_model.ShowSubTypeSeriesTV:   {0.4, 0, 0.1, 0.15, 0.05, 0.3, 0},
	show_model.ShowSubTypeSeriesNet:  {0.4, 0, 0.1, 0.15, 0.05, 0.3, 0},
	show_model.ShowSubTypeMovieNet:   {0.4, 0.3, 0, 0, 0, 0.1, 0.2},
	show_model.ShowSubTypeAmineChina: {0.1, 0.2, 0.1, 0.15, 0.05, 0.4, 0},
}

var today = time.Today()

// 计算骨朵热度
func guduoHotCalc() {
	allHot := make(map[int64]*task.GuduoHotRank)
	data := guduoHotGetInfo()

	totalCommentCount := make([]float64, len(data))
	totalCommentCount2 := make([]float64, len(data))
	totalDanmakuCount := make([]float64, len(data))
	baiduAttention := make([]float64, len(data))
	qihooIndicator := make([]float64, len(data))
	newsCount := make([]float64, len(data))
	wechatArticleCount := make([]float64, len(data))
	baiduIndicator := make([]float64, len(data))
	baiduInfoIndex := make([]float64, len(data))
	baiduIndicatorYesterday := make([]float64, len(data))
	weiboIndicator := make([]float64, len(data))
	weiboIndicatorYesterday := make([]float64, len(data))
	hotWeibo := make([]float64, len(data))
	hot := make([]float64, len(data))
	additions := make([]float64, len(data))

	sss := make([]uint64, 0, 100)
	for i, row := range data {
		sss = append(sss, row.Sid)
		totalCommentCount[i] = row.TotalCommentCount
		totalCommentCount2[i] = 0 // 可能是泡泡评论？现在没有这个东西
		totalDanmakuCount[i] = row.TotalDanmakuCount
		baiduAttention[i] = row.BaiduAttention
		qihooIndicator[i] = row.QihooIndicator
		newsCount[i] = row.BaiduNewsCount
		wechatArticleCount[i] = row.WechatArticleCount
		baiduInfoIndex[i] = 0
		baiduIndicator[i] = row.BaiduIndicator
		baiduIndicatorYesterday[i] = row.BaiduIndicatorYesterday
		weiboIndicator[i] = row.WeiboIndicator
		weiboIndicatorYesterday[i] = row.WeiboIndicatorYesterday
		hotWeibo[i] = row.HotWeibo
		hot[i] = 0
		additions[i] = 0
	}

	// 查看数据
	//fmt.Println("show id", sss)
	//fmt.Println("总评论数", totalCommentCount)
	//fmt.Println("总评论数2-未采集", totalCommentCount2)
	//fmt.Println("总弹幕", totalDanmakuCount)
	//fmt.Println("百度关注人数", baiduAttention)
	//fmt.Println("360指数", qihooIndicator)
	//fmt.Println("新闻数", newsCount)
	//fmt.Println("微信文章数", wechatArticleCount)
	//fmt.Println("百度资讯指数-未采集", baiduInfoIndex)
	//fmt.Println("百度指数", baiduIndicator)
	//fmt.Println("百度指数-昨日", baiduIndicatorYesterday)
	//fmt.Println("微博指数", weiboIndicator)
	//fmt.Println("微博指数-昨日", weiboIndicatorYesterday)
	//fmt.Println("相关微博数", hotWeibo)
	//fmt.Println("hot-未采集", hot)
	//fmt.Println("additions-未采集", additions)
	//return

	totalCommentCount = normalizeAbsolutely(totalCommentCount)
	totalCommentCount2 = normalizeAbsolutely2(totalCommentCount2)
	totalDanmakuCount = normalizeAbsolutely(totalDanmakuCount)
	baiduAttention = normalize(baiduAttention)
	qihooIndicator = normalize(qihooIndicator)
	newsCount = normalize(newsCount)
	wechatArticleCount = normalize(wechatArticleCount)
	baiduInfoIndex = normalize(baiduInfoIndex)
	baiduIndicator = normalize(baiduIndicator)
	baiduIndicatorYesterday = normalize(baiduIndicatorYesterday)
	weiboIndicator = normalize(weiboIndicator)
	weiboIndicatorYesterday = normalize(weiboIndicatorYesterday)
	hotWeibo = normalize(hotWeibo)
	hot = normalize(hot)
	additions = normalize(additions)

	save := make([]*guduo_hot_daily_model.Table, 0, 10000)

	for i, row := range data {
		sid := row.Sid
		if _, ok := allHot[row.Type]; !ok {
			allHot[row.Type] = &task.GuduoHotRank{
				Sid:  row.Sid,
				Type: row.SubType,
				Num:  make([]*task.GuduoHotItem, 0, 300),
			}
		}

		factor, ok := factors[row.SubType]
		if !ok {
			continue
		}

		// 评论
		cmt := .0
		if totalCommentCount2[i] > 0 {
			cmt = totalCommentCount2[i]
		} else {
			cmt = totalCommentCount[i]
		}
		comment := factor[0] * cmt

		// 弹幕
		danmaku := factor[1] * totalDanmakuCount[i]

		// 贴吧关注数
		tieba := factor[2] * baiduAttention[i]

		// 新闻资讯
		media := factor[3] * getNearestMean(qihooIndicator[i], newsCount[i], baiduInfoIndex[i])

		wechat := factor[4] * wechatArticleCount[i]

		// 指数
		idx1 := weiboIndicator[i]
		idx2 := baiduIndicator[i]
		//fmt.Println("----------------")
		//fmt.Println(idx1, idx2)
		//fmt.Println("----------------")
		if weiboIndicator[i] <= 0 {
			idx1 = math.Max(weiboIndicatorYesterday[i], hotWeibo[i])
		}
		if baiduIndicator[i] == 0 {
			idx2 = baiduIndicatorYesterday[i]
		}
		index := factor[5] * getMeanOrEither(idx1, idx2)
		hot_ := factor[6] * hot[i]
		addition := additions[i]

		multi := 1.0
		if comment+danmaku < 0.02 && media+wechat+index > 0.3 {
			multi = 0.5
		}

		//fmt.Println(comment)
		//fmt.Println(danmaku)
		//fmt.Println(tieba)
		//fmt.Println(media)
		//fmt.Println(wechat)
		//fmt.Println(index)
		//fmt.Println(multi)
		//fmt.Println(hot_)
		//fmt.Println(addition)
		//return

		guduoHot := comment + danmaku + tieba + (media+wechat+index)*multi + hot_ + addition
		guduoHot = util.ToFixedFloat(guduoHot, 10)

		// 这是什么意思？
		onBillBoard := 1.0
		if guduoHot > 0 {
			guduoHot = guduoHot*100 + onBillBoard*8
		} else {
			guduoHot = 0
		}

		allHot[row.Type].Num = append(allHot[row.Type].Num, &task.GuduoHotItem{row.SubType, guduoHot})
		//guduo_hot_daily_model.SaveCurHot(guduoHot, today, sid)

		save = append(save, &guduo_hot_daily_model.Table{
			ShowId:    sid,
			Num:       guduoHot,
			DayAt:     today,
		})


	}

	guduo_hot_daily_model.Model().Where("day_at", today).Delete(nil)
	saveFin := make([]*guduo_hot_daily_model.Table, 0, 400)
	for _, row := range save {
		saveFin = append(saveFin, row)
		if len(saveFin) >= 400 {
			guduo_hot_daily_model.Model().Create(&saveFin)
			saveFin = save[:0]
		}
	}

	if len(saveFin) > 0 {
		guduo_hot_daily_model.Model().Create(&saveFin)
	}

	// guduo_hot_daily 存储每天的热度变化，昨日热度有其他任务继续
	//var save []*guduo_hot_rank_model.Table
	//for type_, row := range allHot {
	//	save = append(save, &guduo_hot_rank_model.Table{
	//		ShowId:      row.Sid,
	//		RankType:    model_clean.CycleDaily,
	//		ShowType:    type_,
	//		SubShowType: row.Type,
	//		PlatformId:  0,
	//		Num:         0,
	//		CustomNum:   0,
	//		Rank:        0,
	//		Rise:        0,
	//		DayAt:       0,
	//	})
	//	guduo_hot_rank_model.SaveCurRank(row.Num, type_, model_clean.CycleDaily, JobAt, row.Sid)
	//}
}

func guduoHotGetInfo() []*task.GuduoHotIndicator {
	day := today
	yesterday := day - 86400

	shows := show_model.GetActiveShowsWithType()

	data := make([]*task.GuduoHotIndicator, 0, 10000)

	sids := make([]uint64, 0, len(shows))

	for _, show := range shows {
		sids = append(sids, show.ID)
	}


	TotalCommentCount := comment_count_trend_daily_model.GetCommentCount(sids, []uint{day})
	TotalDanmakuCount := danmaku_count_trend_daily_model.GetDanmakuCount(sids, []uint{day})
	BaiduAttention := attention_trend_daily_model.GetAttention(sids, []uint{day}, constant.PlatformIdBaidu)
	QihooIndicator := indicator_daily_model.GetIndicator(sids, []uint{day}, constant.PlatformIdQihu)
	WechatArticleCount := article_count_trend_daily_model.GetArticleNum(sids, []uint{day}, constant.PlatformIdWeixin)
	NewsCount := news_trend_daily_model.GetNewsCount(sids, []uint{day}, constant.PlatformIdBaidu)
	BaiduNewsCount := news_trend_daily_model.GetNewsCount(sids, []uint{day}, constant.PlatformIdBaidu)
	WeiboIndicator := indicator_daily_model.GetIndicator(sids, []uint{day}, constant.PlatformIdWeibo)
	WeiboIndicatorYesterday := indicator_daily_model.GetIndicator(sids, []uint{yesterday}, constant.PlatformIdWeibo)
	BaiduIndicator := indicator_daily_model.GetIndicator(sids, []uint{day}, constant.PlatformIdBaidu)
	BaiduIndicatorYesterday := indicator_daily_model.GetIndicator(sids, []uint{yesterday}, constant.PlatformIdBaidu)
	HotWeibo := article_count_trend_daily_model.GetArticleNum(sids, []uint{day}, constant.PlatformIdWeibo)

	for _, show := range shows {
		row := &task.GuduoHotIndicator{
			Sid:                     show.ID,
			Type:                    show.ShowType,
			SubType:                 show.SubShowType,
			PlatformIds:             show_model.GetPlatform(show.Platform),
			TotalCommentCount:       0,
			TotalDanmakuCount:       0,
			BaiduAttention:          0,
			QihooIndicator:          0,
			WechatArticleCount:      0,
			NewsCount:               0,
			BaiduNewsCount:          0,
			WeiboIndicator:          0,
			WeiboIndicatorYesterday: 0,
			BaiduIndicator:          0,
			BaiduIndicatorYesterday: 0,
			HotWeibo:                0,
		}

		// 评论数
		for _, item := range TotalCommentCount {
			if item.ShowId == show.ID {
				row.TotalCommentCount = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		//弹幕数
		for _, item := range TotalDanmakuCount {
			if item.ShowId == show.ID {
				row.TotalDanmakuCount = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 百度帖子数
		for _, item := range BaiduAttention {
			if item.ShowId == show.ID {
				row.BaiduAttention = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 360指数
		for _, item := range QihooIndicator {
			if item.ShowId == show.ID {
				row.QihooIndicator = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 微信文章数
		for _, item := range WechatArticleCount {
			if item.ShowId == show.ID {
				row.WechatArticleCount = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 新闻数
		for _, item := range NewsCount {
			if item.ShowId == show.ID {
				row.NewsCount = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 百度新闻数
		for _, item := range BaiduNewsCount {
			if item.ShowId == show.ID {
				row.BaiduNewsCount = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 微博指数
		for _, item := range WeiboIndicator {
			if item.ShowId == show.ID {
				row.WeiboIndicator = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 昨日微博指数
		for _, item := range WeiboIndicatorYesterday {
			if item.ShowId == show.ID {
				row.WeiboIndicatorYesterday = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 百度指数
		for _, item := range BaiduIndicator {
			if item.ShowId == show.ID {
				row.BaiduIndicator = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 昨日百度指数
		for _, item := range BaiduIndicatorYesterday {
			if item.ShowId == show.ID {
				row.BaiduIndicatorYesterday = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}
		// 热门微博数
		for _, item := range HotWeibo {
			if item.ShowId == show.ID {
				row.HotWeibo = util.ToFixedFloat(float64(item.Num), 2)
				break
			}
		}

		data = append(data, row)
	}

	return data
}

func guduoTopK(arr []*task.GuduoHotItem, nc int) []*task.GuduoHotItem {
	maxIdx := len(arr)
	if maxIdx == 1 {
		return arr
	}

	rand.Seed(time2.Now().Unix())
	idx := rand.Intn(maxIdx)
	cnt := arr[idx].Hot

	left := make([]*task.GuduoHotItem, 0, 50)
	mid := make([]*task.GuduoHotItem, 0, 50)
	right := make([]*task.GuduoHotItem, 0, 50)
	for i := 0; i < maxIdx; i++ {
		if arr[i].Hot > cnt {
			left = append(left, arr[i])
		} else if arr[i].Hot == cnt {
			mid = append(mid, arr[i])
		} else {
			right = append(right, arr[i])
		}
	}

	lLen := len(left)
	mLen := len(mid)

	if lLen < nc && mLen+lLen >= nc {
		x := nc - lLen
		mTmp := mid[0:x]

		left = append(left, mTmp...)
		return left
	} else if lLen < nc && mLen+lLen < nc {
		nextCn := nc - (lLen + mLen)
		rTmp := guduoTopK(right, nextCn)

		left = append(left, mid...)

		left = append(left, rTmp...)
		return left
	} else if lLen == nc {
		return left
	} else {
		return guduoTopK(left, nc)
	}
}

var log3 = math.Log(3)
var log4 = math.Log(4)

func normalizeAbsolutely(xs []float64) []float64 {
	for i, x := range xs {
		xs[i] = math.Pow(math.Log(math.Max(1, x))/log3, 2) * .01
	}

	return xs
}
func normalizeAbsolutely2(xs []float64) []float64 {
	for i, x := range xs {
		xs[i] = math.Pow(math.Log(math.Max(1, x))/log4, 2) * .01
	}

	return xs
}

func normalize(xs []float64) []float64 {
	cp := make([]float64, len(xs))
	copy(cp, xs)

	max := .0
	for _, v := range cp {
		if v > max {
			max = v
		}
	}

	max = math.Log(max)

	for k, x := range xs {
		xs[k] = math.Log(math.Max(1, x)) / max
	}

	return xs
}

func getNearestMean(a, b, c float64) float64 {
	mean := (a + b) / 2
	diff := math.Abs(a - b)

	if math.Abs(a-c) < diff {
		diff = math.Abs(a - c)
		mean = (a + c) / 2
	}

	if math.Abs(b-c) < diff {
		mean = (b + c) / 2
	}

	return mean
}

func getMeanOrEither(a, b float64) float64 {
	if a > 0 && b > 0 {
		return (a + b) / 2
	}

	return math.Max(b, math.Max(0, a))
}
