package daily

import (
	"guduo/app/crawler/clean/task"
	"guduo/app/internal/constant"
	"guduo/app/internal/model_clean"
	"guduo/app/internal/model_clean/actor_hot_daily_model"
	"guduo/app/internal/model_clean/actor_hot_rank_model"
	"guduo/app/internal/model_clean/article_count_actor_trend_daily_model"
	"guduo/app/internal/model_clean/attention_actor_trend_daily_model"
	"guduo/app/internal/model_clean/indicator_actor_daily_model"
	"guduo/app/internal/model_scrawler/actor_model"
	"guduo/pkg/time"
	"math"
	"math/rand"
	time2 "time"
)

const (
	weiboScale  float64 = 0.15
	fansScale   float64 = 0.35
	baiduScale  float64 = 0.1
	mediaScale  float64 = 0.15
	wechatScale float64 = 0.25
	baseScale   float64 = 100
)

func GuduoActorHotHandle() {
	guduoActorHot()
}

// 骨朵艺人热力榜
func guduoActorHot() {

	allHot := make([]*task.ActHotItem, 0, 500)
	newHot := make([]*task.ActHotItem, 0, 500)
	curYear := time2.Now().Year()
	data := guduoActorHotGetInfo()

	weiboIndicator := make([]float64, len(data))
	weiboIndicatorYesterday := make([]float64, len(data))
	weiboSpuertopic := make([]float64, len(data))

	weiboIndicatorHour := make([]float64, len(data))
	hotWeiboFeed := make([]float64, len(data))
	hotWeiboComment := make([]float64, len(data))
	hotWeiboForward := make([]float64, len(data))
	hotWeiboLike := make([]float64, len(data))

	baiduIndicator := make([]float64, len(data))
	baiduIndicatorYesterday := make([]float64, len(data))

	baiduIndicatorHour := make([]float64, len(data))
	BaiduInformationIndicator := make([]float64, len(data))

	baiduAttention := make([]float64, len(data))
	qihooIndicator := make([]float64, len(data))

	newsCount := make([]float64, len(data))

	wechatIndex := make([]float64, len(data))
	wechatArticleCount := make([]float64, len(data))
	for i, row := range data {
		weiboIndicator[i] = row.WeiboIndicator
		weiboIndicatorYesterday[i] = row.WeiboIndicatorYesterday
		weiboSpuertopic[i] = 0

		weiboIndicatorHour[i] = 0
		hotWeiboFeed[i] = 0
		hotWeiboComment[i] = 0
		hotWeiboForward[i] = 0
		hotWeiboLike[i] = 0

		baiduIndicator[i] = row.BaiduIndicator
		baiduIndicatorYesterday[i] = row.BaiduIndicatorYesterday

		baiduIndicatorHour[i] = 0
		BaiduInformationIndicator[i] = 0

		baiduAttention[i] = row.BaiduAttention
		qihooIndicator[i] = row.QihooIndicator

		newsCount[i] = 0

		wechatIndex[i] = 0
		wechatArticleCount[i] = row.WechatArticleCount
	}

	weiboIndicator = normalize(weiboIndicator)
	weiboIndicatorYesterday = normalize(weiboIndicatorYesterday)
	weiboSpuertopic = normalize(weiboSpuertopic)

	weiboIndicatorHour = normalize(weiboIndicatorHour)
	hotWeiboFeed = normalize(hotWeiboFeed)
	hotWeiboComment = normalize(hotWeiboComment)
	hotWeiboForward = normalize(hotWeiboForward)
	hotWeiboLike = normalize(hotWeiboLike)

	baiduIndicator = normalize(baiduIndicator)
	baiduIndicatorYesterday = normalize(baiduIndicatorYesterday)

	baiduIndicatorHour = normalize(baiduIndicatorHour)
	BaiduInformationIndicator = normalize(BaiduInformationIndicator)

	baiduAttention = normalize(baiduAttention)
	qihooIndicator = normalize(qihooIndicator)

	newsCount = normalize(newsCount)

	wechatIndex = normalize(wechatIndex)
	wechatArticleCount = normalize(wechatArticleCount)

	for i, row := range data {
		aid := row.Aid

		weibo := weiboScale * weiboIndicator[i]
		if weiboIndicator[i] == 0 {
			curHot := calHotWeibo(hotWeiboFeed[i], hotWeiboComment[i], hotWeiboForward[i], hotWeiboLike[i])

			weibo = math.Max(math.Max(weiboIndicatorYesterday[i], weiboIndicatorHour[i]), curHot)
		}

		fans := getMeanOrEither(weiboSpuertopic[i], baiduAttention[i])

		idx := baiduIndicator[i]
		if idx <= 0 {
			idx = baiduIndicatorHour[i]

			if idx <= 0 {
				idx = baiduIndicatorYesterday[i]
			}
		}
		baidu := baiduScale * idx

		media := mediaScale * getNearestMean(BaiduInformationIndicator[i], newsCount[i], qihooIndicator[i])

		// 微信指数
		wechatIdx := wechatIndex[i]
		if wechatIdx <= 0 {
			wechatIdx = 0.9 * wechatArticleCount[i]
		}
		wechat := wechatScale * wechatIdx

		guduoHot := baseScale * (weibo + fans + baidu + media + wechat)

		allHot = append(allHot, &task.ActHotItem{row.Aid, row.Name, guduoHot})
		isNew := int8(0)
		if int64(curYear)-row.Year <= 25 {
			newHot = append(newHot, &task.ActHotItem{row.Aid, row.Name, guduoHot})
			isNew = 1
		}
		actor_hot_daily_model.SaveCurHot(guduoHot, isNew, JobAt, aid)
	}

	actor_hot_rank_model.SaveCurRank(allHot, 0, model_clean.CycleDaily, JobAt)
	actor_hot_rank_model.SaveCurRank(newHot, 1, model_clean.CycleDaily, JobAt)
}

func guduoActorHotGetInfo() []*task.GuduoActorHotIndicator {
	day := time.Today()
	yesterday := day - 86400

	aids := actor_model.GetActor()

	data := make([]*task.GuduoActorHotIndicator, 0, 100)
	for _, aid := range aids {
		row := &task.GuduoActorHotIndicator{
			Aid:                     aid.Id,
			Name:                    aid.Name,
			Year:                    aid.BirthYear,
			BaiduAttention:          float64(attention_actor_trend_daily_model.GetAttention(aid.Id, []uint{day}, constant.PlatformIdBaidu)),
			QihooIndicator:          float64(indicator_actor_daily_model.GetIndicator(aid.Id, []uint{day}, constant.PlatformIdQihu)),
			WechatArticleCount:      float64(article_count_actor_trend_daily_model.GetArticleNum(aid.Id, []uint{day}, constant.PlatformIdWeixin)),
			WeiboIndicator:          float64(indicator_actor_daily_model.GetIndicator(aid.Id, []uint{day}, constant.PlatformIdWeibo)),
			WeiboIndicatorYesterday: float64(indicator_actor_daily_model.GetIndicator(aid.Id, []uint{yesterday}, constant.PlatformIdWeibo)),
			BaiduIndicator:          float64(indicator_actor_daily_model.GetIndicator(aid.Id, []uint{day}, constant.PlatformIdBaidu)),
			BaiduIndicatorYesterday: float64(indicator_actor_daily_model.GetIndicator(aid.Id, []uint{yesterday}, constant.PlatformIdBaidu)),
		}

		data = append(data, row)
	}

	return data
}

func actorTopK(arr []*task.ActHotItem, nc int) []*task.ActHotItem {
	maxIdx := len(arr)
	if maxIdx == 1 {
		return arr
	}

	rand.Seed(time2.Now().Unix())
	idx := rand.Intn(maxIdx)
	cnt := arr[idx].Hot

	left := make([]*task.ActHotItem, 0, 50)
	mid := make([]*task.ActHotItem, 0, 50)
	right := make([]*task.ActHotItem, 0, 50)
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
		rTmp := actorTopK(right, nextCn)

		left = append(left, mid...)

		left = append(left, rTmp...)
		return left
	} else if lLen == nc {
		return left
	} else {
		return actorTopK(left, nc)
	}
}

func normalize(xs []float64) []float64 {
	max := .0
	ret := make([]float64, len(xs))

	for _, v := range xs {
		if v > max {
			max = v
		}
	}

	if max <= 1 {
		return ret
	}

	max = math.Log(max)
	for i, v := range xs {
		x := 0.5 + math.Log(v)
		if v < 1 {
			x = 0
		}

		ret[i] = x / (0.5 * max)
	}

	return ret
}

func calHotWeibo(feed, comment, forward, like float64) float64 {
	return 0.3*feed + 0.2*comment + 0.2*forward + 0.3*like
}

func getMeanOrEither(a, b float64) float64 {
	if a > 0 && b > 0 {
		return (a + b) / 2
	}
	return math.Max(b, math.Max(0, a))
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
