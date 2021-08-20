package model

import (
	db2 "guduo/pkg/db"
	"strings"
	"time"
)

func StrToTime(s string) uint {
	if strings.Trim(s, " ") == "" {
		return 0
	}

	layout := "2006-01-02"
	if strings.Contains(s, ":") {
		layout = "2006-01-02 15:04:05"
	}


	tt, _ := time.ParseInLocation(layout, s, time.Local)

	ts := uint(tt.Unix())
	if ts >= 2000000000 || ts < 0{
		 return 0
	}

	return ts
}

type ZMigrate struct {
	Table string
	Idx uint64
}

func LastId(name string) uint64 {
	var res ZMigrate
	db := db2.GetCleanMysqlConn()
	db.Model(&ZMigrate{}).Select("idx").Where("table", name).Find(&res)

	return res.Idx
}