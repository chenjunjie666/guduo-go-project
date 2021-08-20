package util

import (
	"sync"
	"time"
)

var uniqueIdLock = &sync.Mutex{}
func UniqueID() uint64 {
	uniqueIdLock.Lock()
	defer uniqueIdLock.Unlock()
	nowNano := time.Now().UnixNano()
	//time.Sleep(time.Millisecond * 10)
	//flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	//id, err := flake.NextID()
	//if err != nil {
	//	return 0
	//}
	return uint64(nowNano)
}
