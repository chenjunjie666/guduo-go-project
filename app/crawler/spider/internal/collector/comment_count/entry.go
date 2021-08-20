package comment_count

import (
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "评论数"

var JobAt = helper2.GetJobAt()

var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}
func Run() {
	jobNum := 6
	wg.Add(jobNum)
	//wg.Add(1)
	// 获取哔哩哔哩的评论数
	go bilibiliHandle()
	//wg.Wait()

	//wg.Add(1)
	// 获取爱奇艺评论数
	go iqiyiHandle()
	//wg.Wait()

	//wg.Add(1)
	// 获取芒果TV评论数
	go mangoHandle()
	//wg.Wait()

	//wg.Add(1)
	// 获取搜狐TV评论数
	go souhuHandle()
	//wg.Wait()

	//wg.Add(1)
	// 获取腾讯视频评论数
	go tencentHandle()
	//wg.Wait()

	//wg.Add(1)
	// 获取优酷评论数
	go youkuHandle()
	//wg.Wait()

	wg.Wait()

	log.Info("评论爬取结束")
}