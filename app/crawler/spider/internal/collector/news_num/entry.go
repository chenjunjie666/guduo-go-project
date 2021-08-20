// 新闻资讯数指标采集器入口文件
package news_num

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "新闻资讯数"

var JobAt = helper2.GetJobAt()
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 1
	wg.Add(jobNum)

	// 百度新闻资讯数获取器
	go baiduHandle()

	wg.Wait()
}
