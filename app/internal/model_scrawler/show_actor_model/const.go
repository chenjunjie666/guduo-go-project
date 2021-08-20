package show_actor_model

// 演员表演位置
const (
	PlayTypeLead      int8 = iota // 领衔主演
	PlayTypeStar                  // 主演
	PlayTypeSupp                  // 配角
	PlayTypeCame                  // 客串
	PlayTypeGuest                 // 常驻嘉宾
	PlayTypeTempGuest             // 暂定嘉宾
	PlayTypeOtherLead             // 其他领衔主演
	PlayTypeOtherStar             // 其他领衔主演
)

var playTypeMap = map[int8]string{
	PlayTypeLead:  "领衔主演",
	PlayTypeStar:  "主演",
	PlayTypeSupp:  "配角",
	PlayTypeCame:  "客串",
	PlayTypeGuest: "常驻嘉宾",
	PlayTypeTempGuest: "暂定嘉宾",
	PlayTypeOtherStar: "其他领衔主演",
}

func GetPlayTypeMap() map[int8]string {
	return playTypeMap
}
