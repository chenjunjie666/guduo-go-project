// 文章数指标采集器入口文件
package article_num_actor

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

var JobAt = helper2.GetJobAt()

const ModName = "艺人文章数"
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 1
	wg.Add(jobNum)
	// 微信文章数获取器
	go wechatHandle()

	wg.Wait()
}
