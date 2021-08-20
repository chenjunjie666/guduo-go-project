package article_num_actor

import (
	"guduo/app/crawler/spider/internal/core"
	"testing"
)

func TestRun(t *testing.T) {
	core.Init()
	Run()
}