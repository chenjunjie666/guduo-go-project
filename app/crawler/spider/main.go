package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/config"
	"guduo/app/crawler/spider/internal/core"
	"os"
	"strings"
	"time"
)

var WorkName = ""
var WorkStartTime = time.Now().Format("2006-01-02_15_04")

func main() {
	runAll()
	//start()
	//runOne()
}

func start() {
	log.Info("======================== START ========================")
	fname := os.Args[0]
	args := os.Args[1:len(os.Args)]
	if len(args) == 0 {
		fmt.Println(fmt.Sprintf("输入 \"%s help\" 查看帮助", fname))
		return
	}

	if len(args) == 1 && args[0] == "help" {
		helpStr :=
			fmt.Sprintf(`
+++++++++++++++++++
用法： %s job1 [job2 job3 ...]
[xxx]表示可选参数
+++++++++++++++++++
可用的job：
all                全部
article_content    微博文章内容
article_num        微信文章数
article_num_actor  微信文章数-艺人
attention          百度贴吧关注数&发帖数
attention_actor    百度贴吧关注数&发帖数-艺人
base_info          主创&公司&演员
comment_count      视频评论数
danmaku            弹幕内容&弹幕数 #注意 当danmaku为第一个参数，且第二个参数为 bilibili tencent mango souhu iqiyi youku 时，可以只爬对应平台的弹幕
fans               微博粉丝数
hot                热度趋势
indicator          百度&搜狗&微博指数
indicator_actor    指数-艺人
introduction       剧情简介
length             视频单集时长
news_num           百度新闻数
play_count         播放量
rating_num         豆瓣评分
release_time       上线时间
short_comment      豆瓣短评数
`, fname)
		fmt.Print(helpStr)
		return
	}
	fs := make([]config.JobFunc, 0, 5)

	if args[0] == "all" {
		for _, v := range config.JobMap {
			fs = append(fs, v)
		}
	} else {
		jobs := config.JobMap
		for _, arg := range args {
			if jobs[arg] == nil {
				log.Warning(fmt.Sprintf("任务：%s 不存在，请确认任务名是否正确，请输入 \"%s help\" 查看帮助", arg, fname))
				continue
			}
			fs = append(fs, jobs[arg])
		}
	}

	go func() {
		time.Sleep(time.Hour * 23)
		jobName := strings.Join(os.Args, " ")
		log.Fatal("任务名", jobName, "超过最大执行时间-23小时")
		panic("任务名" + jobName + "超过最大执行时间-23小时")
	}()

	core.Init()
	for _, f := range fs {
		f()
	}

	log.Info("======================== DONE ========================")
}


func runAll(){
	core.Init()
	//config.JobMap["news_num"]()
	fs := make([]config.JobFunc, 0, 5)
	for _, v := range config.JobMap {
		fs = append(fs, v)
	}
	for _, f := range fs {
		f()
	}
}

func runOne(){
	core.Init()
	config.JobMap["danmaku"]()
}