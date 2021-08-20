package api

import (
	"guduo/app/cms/test"
	"testing"
)

func TestShowList(t *testing.T) {
	data := map[string]string{
		"show_type": "0",
	}


	req := test.PostRequest("/show/list?debug=1", data)

	test.PrintRespJson(req)
}
