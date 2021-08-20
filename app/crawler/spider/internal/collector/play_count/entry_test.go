package play_count

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestRun(t *testing.T) {
	core.Init()
	Run()
}
