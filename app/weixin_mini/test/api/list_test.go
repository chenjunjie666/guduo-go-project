package api

import (
	"guduo/app/weixin_mini/test"
	"testing"
)

func TestShowList(t *testing.T) {
	data := map[string]string{
		"day_at": "1621699200",
		"list_type": "0",
		"type": "0",
		"sub_type": "-1",
		"platform_id": "0",
	}

	req := test.GetRequest("/show/home/list", data)

	test.PrintRespJson(req)
}

func TestActorList(t *testing.T) {
	data := map[string]string{
		//"day_at": "1619366400",
		"day_at": "1620835200",
		"list_type": "3",
		"rank_type": "1",
		"play_type": "0",
	}

	req := test.GetRequest("/actor/home/list", data)

	test.PrintRespJson(req)
}


func TestSearch(t *testing.T) {
	data := map[string]string{
		"keyword": "武当一剑" ,
	}

	req := test.GetRequest("/show/home/search", data)

	test.PrintRespJson(req)
}




func TestSearchHot(t *testing.T) {
	data := map[string]string{
		"type": "0",
	}

	req := test.GetRequest("/show/home/hot_search", data)

	test.PrintRespJson(req)
}

