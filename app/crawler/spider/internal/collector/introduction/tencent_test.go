package introduction

import (
	"guduo/app/crawler/spider/internal/core"
	"guduo/app/crawler/spider/internal/storage"
	"testing"
)

func TestTencentIntroduce(t *testing.T) {
	core.Init()

	wg.Add(1)
	url := storage.Tencent.GetDetailUrl()[0]
	tencentIntroduction(url.Url, url.ShowId)
}
