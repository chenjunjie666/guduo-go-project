package base_info

import (
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "基本信息"

var JobAt = helper2.GetJobAt()

var wg = &sync.WaitGroup{}

// todo 综艺嘉宾类型
// 对于不同类型，是否有不同的处理方式，等待测试ing
func Run() {
	jobNum := 5
	wg.Add(jobNum)

	// 爱奇艺基本信息
	go iqiyiHandle()

	// 芒果TV基本信息
	go mangoHandle()

	// 搜狐基本信息
	go souhuHandle()

	// 腾讯基本信息
	go tencentHandle()

	// 优酷基本信息
	go youkuHandle()

	wg.Wait()
}
