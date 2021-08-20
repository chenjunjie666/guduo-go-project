package util

import (
	"strconv"
	"strings"
)

func EscapeDotInt(s string) (int64, error) {
	sp := strings.Split(s, ",")
	numStr := strings.Join(sp, "")
	num, e := strconv.ParseInt(numStr, 10, 64)
	return num, e
}
