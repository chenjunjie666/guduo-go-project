package comment_count

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

// todo 待验证正确性
func TestRun(t *testing.T) {
	core.Init()
	Run()
}