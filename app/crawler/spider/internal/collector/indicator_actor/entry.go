package indicator_actor

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "演员指数"

var JobAt = helper2.GetJobAt()
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}

func Run() {
	jobNum := 3
	wg.Add(jobNum)

	// 百度指数
	go baiduHandle()

	// 微博指数
	go weiboHandle()

	// 360指数
	go qihooHandle()

	wg.Wait()
}
