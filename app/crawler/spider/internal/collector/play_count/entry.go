// 播放量指标采集器入口文件
package play_count

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "视频播放量"

var JobAt = helper2.GetJobAt()
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 2
	wg.Add(jobNum)

	// 芒果TV播放量
	go mangoHandle()

	// 腾讯视频播放量
	go tencentHandle()

	wg.Wait()
}
