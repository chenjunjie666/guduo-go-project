package show_model

const (
	ShowTypeUnknown     int64 = iota - 1 // 未知类型
	ShowTypeSeries                       // 剧集
	ShowTypeVariety                      // 综艺
	ShowTypeMovie                        // 电影
	ShowTypeMicro                        // 竖屏微剧
	ShowTypeAmine                        // 动漫
	ShowTypeDocumentary                  // 纪录片
)

const (
	// 0 开始是剧集分类
	ShowSubTypeUnknown       int64 = iota - 1 // 未知子类型
	ShowSubTypeSeriesNet                      // 网络剧
	ShowSubTypeSeriesTV                       // 电视剧
	ShowSubTypeSeriesAmerica                  // 美剧
	ShowSubTypeSeriesJP                       // 日剧
	ShowSubTypeSeriesKR                       // 韩剧
	_
	_
	_
	_
	_ // 2-9做剧集子分类预留

	// 10 开始是综艺分类
	ShowSubTypeVarietyNet // 网综 - 网络综艺
	ShowSubTypeVarietyTV  // 电视综艺
	_
	_
	_
	_
	_
	_
	_
	_ // 12-19做剧综子分类预留

	// 20 开始是电影分类
	ShowSubTypeMovieNet    // 网络电影
	ShowSubTypeMovieCinema // 院线网播
	_
	_
	_
	_
	_
	_
	_
	_ // 22-29 做电影子分类预留

	// 30 开始是微剧子分类
	_
	_
	_
	_
	_
	_
	_
	_
	_
	_ // 30-39 做微剧子分类预留

	// 40 开始是动漫分类
	ShowSubTypeAmineChina           // 国漫
	ShowSubTypeAmineJP              // 日本动画
	ShowSubTypeAmineForeignKid      // 外国儿童动画
	ShowSubTypeAmineMovieForeignKid // 外国儿童动画电影
	ShowSubTypeAmineMovieKid        // 儿童动画电影
	_
	_
	_
	_
	_ // 40-49 做动漫子分类预留
	// 51 开始是纪录片
	ShowTypeSubDocumentary // 纪录片（目前没有别的分类

)

const (
	ShowStatReject   int64 = iota - 1 // 审核不通过 -1
	ShowStatPending                   // 待审核 0
	ShowStatStandard                  // 正常状态 1
)

const (
	ShowOff = iota // 剧集不显示在小程序端
	ShowOn         // 剧集显示在小程序端
)

const (
	ShowSearchHot    = iota // 剧集是热搜
	ShowSearchNotHot        // 剧集不是热搜
)

const (
	ShowPlayingStatPlaying int = iota // 在播
	ShowPlayingStatWaiting            // 待播
	ShowPlayingStatOff                // 下架
	ShowPlayingStatInvalid            // 无效
)
