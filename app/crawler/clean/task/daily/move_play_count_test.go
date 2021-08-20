package daily

import (
	"guduo/app/crawler/clean/internal/core"
	"testing"
)

func TestMoivePlayCount(t *testing.T) {
	core.Init()
	moviePlayCountHandle()
}
