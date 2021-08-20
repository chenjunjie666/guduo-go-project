package daily

import (
	"guduo/app/crawler/clean/internal/core"
	"testing"
)

func TestAvgGuduoHotHandle(t *testing.T) {
	core.Init()
	AvgGuduoHotHandle()
}