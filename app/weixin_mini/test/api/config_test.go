package api

import (
	"guduo/app/weixin_mini/test"
	"testing"
)

func TestConfig(t *testing.T) {
	req := test.GetRequest("/config", nil)

	test.PrintRespJson(req)
}

func TestDateConfig(t *testing.T) {

	data := map[string]string{
		"actor": "1",
	}
	req := test.GetRequest("/date_config", data)

	test.PrintRespJson(req)
}
