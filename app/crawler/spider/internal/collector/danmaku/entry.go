// 弹幕内容指标采集器入口文件
package danmaku

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"os"
	"sync"
)

const ModName = "弹幕内容"

var JobAt = helper2.GetJobAt()
var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}

func Run() {
	jobNum := 6
	//jobNum := 1

	args := ""
	if len(os.Args) == 3 {
		args = os.Args[2]
	}

	switch args {
	case "bilibili":
		wg.Add(1)
		go bilibiliHandle()
	case "tencent":
		wg.Add(1)
		go tencentHandle()
	case "mango":
		wg.Add(1)
		go mangoHandle()
	case "souhu":
		wg.Add(1)
		go souhuHandle()
	case "iqiyi":
		wg.Add(1)
		go iqiyiHandle()
	case "youku":
		wg.Add(1)
		go youkuHandle()
	default:
		wg.Add(jobNum)

		go bilibiliHandle()
		go tencentHandle()
		go mangoHandle()
		go souhuHandle()
		go iqiyiHandle()
		go youkuHandle()
	}

	// 弹幕内容获取器

	wg.Wait()
	fmt.Println("弹幕爬取结束")
}
