package hot

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "热度"

var JobAt = helper2.GetJobAt()
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 2
	wg.Add(jobNum)

	// 爱奇艺热度趋势
	go iqiyiHandle()

	// 优酷热度
	go youkuHandle()

	wg.Wait()
}
