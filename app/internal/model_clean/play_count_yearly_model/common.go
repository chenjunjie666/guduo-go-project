package play_count_yearly_model

import (
	"fmt"
	"gorm.io/gorm"
	"guduo/pkg/db"
	"strconv"
	"time"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前播放量
func SaveCurPlayCount(num int64, da uint, sid uint64) {
	yaStr := time.Unix(int64(da), 0).Format("2006")
	yaInt, _ := strconv.Atoi(yaStr)
	ya :=time.Date(yaInt, 0, 0, 0, 0, 0, 0, time.Local).Unix()

	var cnt int64
	Model().Where("show_id = ? and day_at = ?", sid, da).Count(&cnt)
	if cnt == 0 {
		row := &Table{
			ShowId:     sid,
			Num:        num,
			DayAt:      uint(ya),
		}
		r := Model().Create(&row)
		if r.Error != nil {
			sql := fmt.Sprintf("update `play_count_yearly` set num = num + %d where show_id=%d and day_at=%d", num, sid, da)
			db.GetCleanMysqlConn().Exec(sql)
		}
	}else{
		sql := fmt.Sprintf("update `play_count_yearly` set num = num + %d where show_id=%d and day_at=%d", num, sid, da)
		db.GetCleanMysqlConn().Exec(sql)
	}
}