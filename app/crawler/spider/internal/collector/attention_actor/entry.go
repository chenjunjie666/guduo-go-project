// 贴数指标采集器入口文件
package attention_actor

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "艺人贴吧关注度"

var wg = &sync.WaitGroup{}
var ch = core.NewJobQueue(40)
var JobAt = helper2.GetJobAt()

func Run() {
	jobNum := 1
	wg.Add(jobNum)
	// 百度贴吧关注度获取器
	go baiduHandle()

	wg.Wait()
}
