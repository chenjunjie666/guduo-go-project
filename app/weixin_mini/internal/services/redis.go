// 由于都在一个db下操作，所以
// 规则为：
// 类型_sid_pid
// 如果该项目没有pid或者表示全平台，pid应为0
package services

const (
	prefixRank = "rank"
)

func GetRank(key string) {
	//rdb := db.GetRedisConn()
}
