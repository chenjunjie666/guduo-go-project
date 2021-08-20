package controller

import (
	"fmt"
	"testing"
)


func TestDeleteProxy(t *testing.T) {
	ipMap["http://123.123.123.123"] = 0
	ipMap["http://123.123.123.124"] = 0
	ipMap["http://123.123.123.125"] = 0
	ipMap["http://123.123.123.126"] = 0
	ipMap["http://123.123.123.127"] = 0

	row := &proxyUrl{"http://123.123.123.123", int64(1), 0, 0}
	pool = append(pool, row)
	row = &proxyUrl{"http://123.123.123.124", int64(1), 0, 0}
	pool = append(pool, row)
	row = &proxyUrl{"http://123.123.123.125", int64(1), 0, 0}
	pool = append(pool, row)
	row = &proxyUrl{"http://123.123.123.126", int64(1), 0, 0}
	pool = append(pool, row)
	row = &proxyUrl{"http://123.123.123.127", int64(1), 0, 0}
	pool = append(pool, row)

	deleteUrl("123.123.123.125")

	fmt.Println(pool)
	fmt.Println(ipMap)
}
