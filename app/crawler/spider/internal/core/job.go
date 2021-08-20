package core

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

var startAt = time.Now().Unix()
func NewJobQueue(num int) *job {
	instance := &job{
		cnt:      0,
		complete: 0,
		ch:       make(chan bool, num),
	}
	return instance
}

type job struct {
	cnt int64
	complete int64
	ch chan bool
}

func (j *job) PushJob() {
	log.Info(">>>>>>>>>>>>>>>>>>当前等待任务数", j.cnt)
	j.ch <- true
	j.cnt++
}

func (j *job) PopJob() {
	<- j.ch
	j.cnt--
	j.complete++
	if j.complete & 15 == 0 {
		// 每完成15个任务手动执行一次垃圾回收
		log.Info("开始强制执行垃圾回收")
		log.Info("垃圾回收执行完毕")
		runtime.GC()
	}
	now := time.Now().Unix()
	sec := now - startAt
	h := sec / 3600
	m := (sec % 3600) / 60
	s := (sec % 3600) % 60

	tStr := ""
	if h > 0 {
		tStr += fmt.Sprintf("%d小时", h)
	}
	if m > 0 || tStr != ""{
		tStr += fmt.Sprintf("%d分钟", m)
	}
	if s > 0 {
		tStr += fmt.Sprintf("%d秒", s)
	}

	log.Info(">>>>>>>>>>>>>>>>>>弹出一个任务，当前等待任务数", j.cnt, " 总计完成任务数", j.complete, "当前耗时：", tStr)
}
