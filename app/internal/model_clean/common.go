package model_clean

const (
	CycleDaily   int8 = iota + 1 // 日榜
	CycleWeekly                  // 周榜
	CycleMonthly                 //月榜
	CycleYearly                  //年榜
	CycleAll                     //总榜
)
