// 片长指标采集器入口文件
package length

import (
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "片长"

var JobAt = helper2.GetJobAt()
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 1
	wg.Add(jobNum)

	// 豆瓣单集片长获取器
	go doubanHandle()

	wg.Wait()
}
