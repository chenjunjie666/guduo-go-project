package api

import (
	"guduo/app/cms/test"
	"testing"
)

func TestDetail(t *testing.T) {
	data := map[string]string{
		"show_id": "1",
	}

	req := test.PostRequest("/show/detail?debug=1", data)

	test.PrintRespJson(req)
}


func TestConfig(t *testing.T) {

	req := test.PostRequest("/show/detail/config?debug=1", nil)

	test.PrintRespJson(req)
}