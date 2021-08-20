package article_content

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"sync"
)

const ModName = "相关文章"

var JobAt = helper2.GetJobAt()

var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}

func Run() {
	jobNum := 1
	wg.Add(jobNum)

	// 爬取微博文章
	go weiboHandle()

	wg.Wait()
}