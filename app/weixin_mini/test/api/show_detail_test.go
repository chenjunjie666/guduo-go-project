package api

import (
	"guduo/app/weixin_mini/test"
	"testing"
)

func TestBaseInfoApi(t *testing.T) {
	data := map[string]string{
		"show_id": "19202",
	}

	req := test.GetRequest("/show/detail/base_info", data)

	test.PrintRespJson(req)
}

func TestNetHot(t *testing.T) {
	data := map[string]string{
		"show_id": "22237",
	}

	req := test.GetRequest("/show/detail/net_hot", data)

	test.PrintRespJson(req)
}

func TestDouban(t *testing.T) {
	data := map[string]string{
		"show_id": "22237",
	}

	req := test.GetRequest("/show/detail/douban", data)

	test.PrintRespJson(req)
}

func TestDanmakuComment(t *testing.T) {
	data := map[string]string{
		"show_id": "22347",
	}

	req := test.GetRequest("/show/detail/danmaku_comment", data)

	test.PrintRespJson(req)
}

func TestWordCloud(t *testing.T) {
	data := map[string]string{
		"show_id": "22347",
	}

	req := test.GetRequest("/show/detail/word_cloud", data)

	test.PrintRespJson(req)
}

func TestWeibo(t *testing.T) {
	data := map[string]string{
		"show_id": "19202",
	}

	req := test.GetRequest("/show/detail/weibo", data)

	test.PrintRespJson(req)
}

func TestWeiboHot(t *testing.T) {
	data := map[string]string{
		"show_id": "1",
	}

	req := test.GetRequest("/show/detail/hot_weibo", data)

	test.PrintRespJson(req)
}

func TestAnalysis(t *testing.T) {
	data := map[string]string{
		"show_id": "1",
	}

	req := test.GetRequest("/show/detail/analysis", data)

	test.PrintRespJson(req)
}