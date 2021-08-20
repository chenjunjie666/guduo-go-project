package task

import (
	"math/rand"
	time2 "time"
)

type GuduoHotIndicator struct {
	Sid                     uint64
	Type                    int64
	SubType                 int64
	PlatformIds             []int
	TotalCommentCount       float64
	TotalDanmakuCount       float64
	BaiduAttention          float64
	QihooIndicator          float64
	WechatArticleCount      float64
	NewsCount               float64 // ???
	BaiduNewsCount          float64
	WeiboIndicator          float64
	WeiboIndicatorYesterday float64
	BaiduIndicator          float64
	BaiduIndicatorYesterday float64
	HotWeibo                float64
}

type GuduoActorHotIndicator struct {
	Aid                     uint64
	Name                    string
	Year                    int64
	BaiduAttention          float64
	QihooIndicator          float64
	WechatArticleCount      float64
	WeiboIndicator          float64
	WeiboIndicatorYesterday float64
	BaiduIndicator          float64
	BaiduIndicatorYesterday float64
}

type GuduoActorDomiIndicator struct {
	Aid      uint64
	Type     int8
	Name     string
	ActorHot float64
	GuduoHot []float64
}

type GuduoHotRank struct {
	Sid  uint64
	Type int64
	Num  []*GuduoHotItem
}

type GuduoHotItem struct {
	SubType int64
	Hot     float64
}

type ActHotItem struct {
	Aid       uint64
	ActorName string
	Hot       float64
}

func GuduoTopK(arr []*GuduoHotItem, nc int) []*GuduoHotItem {
	maxIdx := len(arr)
	if maxIdx == 1 {
		return arr
	}

	rand.Seed(time2.Now().Unix())
	idx := rand.Intn(maxIdx)
	cnt := arr[idx].Hot

	left := make([]*GuduoHotItem, 0, 50)
	mid := make([]*GuduoHotItem, 0, 50)
	right := make([]*GuduoHotItem, 0, 50)
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
		rTmp := GuduoTopK(right, nextCn)

		left = append(left, mid...)

		left = append(left, rTmp...)
		return left
	} else if lLen == nc {
		return left
	} else {
		return GuduoTopK(left, nc)
	}
}

func ActorTopK(arr []*ActHotItem, nc int) []*ActHotItem {
	maxIdx := len(arr)
	if maxIdx == 1 {
		return arr
	}

	rand.Seed(time2.Now().Unix())
	idx := rand.Intn(maxIdx)
	cnt := arr[idx].Hot

	left := make([]*ActHotItem, 0, 50)
	mid := make([]*ActHotItem, 0, 50)
	right := make([]*ActHotItem, 0, 50)
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
		rTmp := ActorTopK(right, nextCn)

		left = append(left, mid...)

		left = append(left, rTmp...)
		return left
	} else if lLen == nc {
		return left
	} else {
		return ActorTopK(left, nc)
	}
}
