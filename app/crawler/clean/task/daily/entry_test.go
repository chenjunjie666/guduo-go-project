package daily

import (
	"guduo/app/crawler/clean/internal/core"
	"testing"
)

func TestRun(t *testing.T) {
	core.Init()
	Run()
}