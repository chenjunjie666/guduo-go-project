package danmaku_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/constant"
	"guduo/pkg/db"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

func GetDanmaku(sid uint64, pid ...uint) []string {
	var dmkRes []*Table

	mm := Model().Select("content")
	if len(pid) > 0 {
		mm.Where("platform_id IN ?", pid)
	}
	mm.Where("show_id = ?", sid).Find(&dmkRes)

	dmk := make([]string, len(dmkRes))

	for k, row := range dmkRes {
		dmk[k] = row.Content
	}

	return dmk
}

// 保存弹幕
func SaveDanmaku(cts []string, jobAt uint, showId, platformId uint64, origin... []string) int64 {
	Model().Where("job_at = ? and show_id = ? and platform_id = ?", jobAt, showId, platformId).
		Delete(nil)

	// 芒果tv无法筛选时间，直接走暴力去重
	if platformId == constant.PlatformIdMango || platformId == constant.PlatformIdTencent {
		ids := doRepeat(origin[0], showId, platformId)

		var d []*Table

		for _, idx := range ids {
			row := &Table{
				ShowId: showId,
				PlatformId: platformId,
				Content: cts[idx],
				ContentId: origin[0][idx],
				JobAt: jobAt,
			}

			d = append(d, row)

			if len(d) >= 400 {
				Model().Create(&d)
				d = d[:0]
			}
		}

		if len(d) > 0 {
			Model().Create(&d)
		}

		return int64(len(ids))
	}

	var d []*Table

	for _, content := range cts {
		row := &Table{
			ShowId: showId,
			PlatformId: platformId,
			Content: content,
			JobAt: jobAt,
		}

		d = append(d, row)

		if len(d) >= 400 {
			Model().Create(&d)
			d = d[:0]
		}
	}

	if len(d) > 0 {
		Model().Create(&d)
	}
	return int64(len(cts))
}

func doRepeat(cmtId []string, sid, pid uint64) []int {
	var res []Table

	//这个whereIn太大了，需要处理一下
	rows := chunkBy(cmtId, 10000)
	for _, row := range rows {
		var chunk []Table
		Model().Select("content_id").
			Where("show_id = ? and platform_id = ?", sid, pid).
			Where("content_id IN ?", row).
			Find(&chunk)

		res = append(res, chunk...)
	}

	ret := make([]int, 0, 1000)
	for k, v := range cmtId {
		repeat := false
		for _, row := range res {
			if v == row.ContentId {
				repeat = true
				break
			}
		}

		if !repeat {
			ret = append(ret, k)
		}
	}

	return ret
}


func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}