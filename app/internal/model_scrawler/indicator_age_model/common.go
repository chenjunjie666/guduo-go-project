package indicator_age_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_clean/indicator_age_daily_model"
	"guduo/pkg/db"
	"strconv"
	"strings"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

// 保存总指数
func SaveIndicatorAge(num map[string]float64, jobAt uint, showId, platformId uint64) {
	var cnt int64
	Model().Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Count(&cnt)

	if cnt > 0{
		return
	}

	var d []*Table
	for a, rating := range num {
		tmp := strings.Split(a, "-")
		// 统计某某+岁的数据
		if strings.Contains(a, "+"){
			tmp_ := strings.Split(a, "-")
			tmp = []string{tmp_[0], "120"}
		}
		if len(tmp) == 2 {
			f, e1 := strconv.ParseInt(tmp[0], 10, 64)
			t, e2 := strconv.ParseInt(tmp[1], 10, 64)
			if e1 != nil || e2 != nil {
				continue
			}
			row := &Table{
				ShowId: showId,
				PlatformId: platformId,
				AgeFrom: f,
				AgeTo: t,
				Rating: rating,
				JobAt: jobAt,
			}
			d = append(d, row)
		}
	}
	Model().Create(&d)

	indicator_age_daily_model.SaveCurAge(num, jobAt, showId, platformId)
}