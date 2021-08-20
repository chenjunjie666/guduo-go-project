package main

import (
	"fmt"
	"guduo/app/crawler/clean/internal/config"
	"guduo/app/crawler/clean/internal/core"
	"os"
)

func main() {
	core.Init()

	//config.JobMap["guduo_hot"]()
	//return


	fname := os.Args[0]
	args := os.Args[1: len(os.Args)]
	if len(args) == 0{
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
daily  每日任务
weekly  每日任务
monthly  每日任务
guduo_hot 计算当日骨朵热度
`, fname)
		fmt.Print(helpStr)
		return
	}

	fs := make([]config.JobFunc, len(args))
	jobs := config.JobMap
	for k, arg := range args {
		if jobs[arg] == nil {
			fmt.Println(fmt.Sprintf("任务：%s 不存在，请确认任务名是否正确，请输入 \"%s help\" 查看帮助", arg, fname))
			return
		}

		fs[k] = jobs[arg]
	}

	core.Init()
	for _, f := range fs {
		f()
	}
}
