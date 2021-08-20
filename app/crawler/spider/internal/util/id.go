package util

import (
	"fmt"
	"time"
)

func GetWorkerId(prefix string) string {
	uniqueWorkerId := fmt.Sprintf("%s__%d", prefix, time.Now().UnixNano())
	return uniqueWorkerId
}
