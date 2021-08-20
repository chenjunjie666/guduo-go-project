package attention_actor

import (
	"fmt"
	"guduo/app/crawler/spider/internal/core"
	"guduo/pkg/util"
	"testing"
)

func TestBaiduAttention(t *testing.T) {
	core.Init()
	wg.Add(1)
	//url := storage.Baidu.GetActorDetailUrl()[0]
	url := fmt.Sprintf("https://tieba.baidu.com/f?ie=utf-8&kw=%s&fr=search", util.UrlEncode("景甜"))
	baiduAttention(url, 0)
}