package introduction

import (
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "剧情介绍"

var JobAt = helper2.GetJobAt()
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 1
	wg.Add(jobNum)
	// 腾讯获取剧情介绍
	go tencentHandle()

	wg.Wait()
}
