package daily

import (
	"guduo/app/crawler/clean/internal/core"
	"testing"
)

func TestYearTotalPlayCountHandle(tt *testing.T) {
	core.Init()

	YearTotalPlayCountHandle()
}