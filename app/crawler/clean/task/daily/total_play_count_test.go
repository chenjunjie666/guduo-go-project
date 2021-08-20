package daily

import (
	"guduo/app/crawler/clean/internal/core"
	"testing"
)

func TestTotalPlayCountTest(t *testing.T) {
	core.Init()


	calcTotalPlayCount()
}
