// 文章数指标采集器入口文件
package article_num

import (
	"guduo/app/crawler/spider/internal/core"
	helper2 "guduo/app/internal/helper"
	"regexp"
	"strconv"
	"sync"
)

var JobAt = helper2.GetJobAt()

const ModName = "文章数"

var ch = core.NewJobQueue(40)
var wg = &sync.WaitGroup{}

func Run() {
	jobNum := 1
	wg.Add(jobNum)
	// 微信文章数获取器
	go wechatHandle()
	// 微博文章数获取器 废弃，相关微博数在 article_content 下获取
	//go weiboHandle()

	wg.Wait()
}

func filterArticleNum(s string) int64 {
	reg := regexp.MustCompile(`\d+`)
	result := reg.FindAllString(s, -1)
	str := ""
	for _, text := range result {
		str += text
	}
	anInt, err := strconv.ParseInt(str, 10, 64)

	if err != nil {
		return 0
	}

	return anInt
}