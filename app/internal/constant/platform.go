package constant

// 平台的ID
const (
	PlatformIdTencent    = 1   // 腾讯
	PlatformIdIqiyi      = 2   // 爱奇艺
	PlatformIdYouku      = 3   // 优酷
	PlatformIdSouhu      = 5   // 搜狐
	PlatformIdMango      = 7   // 芒果
	PlatformIdDouban     = 12  // 豆瓣
	PlatformIdQihu       = 14  // 360
	PlatformIdWeixin     = 16  // 微信
	PlatformIdWeibo      = 17  // 微博
	PlatformIdBilibili   = 24  // bilibili
	PlatformIdBaidu      = 27  // 百度
	PlatformIdSogou      = 100 // 搜狗
	PlatformIdTikTalk    = 110 // 抖音
	PlatformIdKuaishou   = 111 // 快手
	PlatformIdTxMicro    = 112 // 腾讯微视
	PlatformIdKuaidianTV = 113 // 快点TV
	PlatformIdFanYue     = 114 // 番乐
)

var platformMap = map[int]string{
	PlatformIdTencent:    "腾讯",
	PlatformIdYouku:      "优酷",
	PlatformIdMango:      "芒果",
	PlatformIdIqiyi:      "爱奇艺",
	PlatformIdBilibili:   "bilibili",
	PlatformIdWeibo:      "微博",
	PlatformIdDouban:     "豆瓣",
	PlatformIdBaidu:      "百度",
	PlatformIdQihu:       "360",
	PlatformIdSouhu:      "搜狐",
	PlatformIdSogou:      "搜狗",
	PlatformIdWeixin:     "微信",

	PlatformIdTikTalk:    "抖音",
	PlatformIdKuaishou:   "快手",
	PlatformIdTxMicro:    "腾讯微视",
	PlatformIdKuaidianTV: "快点TV",
	PlatformIdFanYue:     "番乐",
}

func GetPlatformMap() map[int]string {
	return platformMap
}

var videoPlatformMap = map[int]string{
	PlatformIdTencent:  "腾讯",
	PlatformIdYouku:    "优酷",
	PlatformIdMango:    "芒果",
	PlatformIdIqiyi:    "爱奇艺",
	PlatformIdSouhu:    "搜狐",
	PlatformIdBilibili: "bilibili",

	// 短视频平台
	PlatformIdTikTalk:    "抖音",
	PlatformIdKuaishou:   "快手",
	PlatformIdTxMicro:    "腾讯微视",
	PlatformIdKuaidianTV: "快点TV",
	PlatformIdFanYue:     "番乐",
}

func GetVideoPlatformMap() map[int]string {
	return videoPlatformMap
}
