package indicator_age_daily_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/helper"
	"guduo/pkg/db"
	"strconv"
	"strings"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCleanMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存当日当前年龄分布
func SaveCurAge(num map[string]float64, ja uint, sid, pid uint64) {
	da := helper.JobAt2DayAt(ja)
	Model().Where("show_id = ? and platform_id = ? and day_at = ?", sid, pid, da).Delete(nil)

	var d []*Table
	for a, rating := range num {
		tmp := strings.Split(a, "-")
		if len(tmp) == 2 {
			f, e1 := strconv.ParseInt(tmp[0], 10, 64)
			t, e2 := strconv.ParseInt(tmp[1], 10, 64)
			if e1 != nil || e2 != nil {
				continue
			}
			row := &Table{
				ShowId: sid,
				PlatformId: pid,
				AgeFrom: f,
				AgeTo: t,
				Rating: rating,
				DayAt: da,
			}
			d = append(d, row)
		}
	}
	Model().Create(&d)
}