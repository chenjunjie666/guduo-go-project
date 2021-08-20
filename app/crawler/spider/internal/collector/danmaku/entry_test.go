package danmaku

import (
	"guduo/app/crawler/spider/internal/core"
	_ "net/http/pprof"
	"testing"
)

func TestMango(t *testing.T) {
	//go http.ListenAndServe("0.0.0.0:8888", nil)
	core.Init()
	ch = core.NewJobQueue(40)
	//wg.Add(1)
	//go mangoHandle()



	jobNum := 6
	wg.Add(jobNum)

	go bilibiliHandle()
	go tencentHandle()
	go mangoHandle()
	go souhuHandle()
	go iqiyiHandle()
	go youkuHandle()

	wg.Wait()
}
