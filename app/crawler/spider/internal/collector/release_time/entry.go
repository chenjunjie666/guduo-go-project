// 上线时间指标采集器入口文件
package release_time

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "视频上线时间"

var JobAt = helper2.GetJobAt()
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 1
	wg.Add(jobNum)

	// 百度百科上线时间
	go baiduHandle()

	wg.Wait()
}
