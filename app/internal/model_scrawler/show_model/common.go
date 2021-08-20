package show_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/test"
	"guduo/pkg/db"
	"guduo/pkg/model"
	"guduo/pkg/time"
)

// 直接获取当前模型的对象
var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

func GetBaseInfo(sid uint64) *ShowBaseInfo {
	var bInfo *ShowBaseInfo

	r := Model().Select("id", "poster",
		"name", "show_type", "sub_show_type", "platform", "tag",
		"staff", "director", "introduction",
		"length", "release_at", "total_episode",
		"show_status").
		Where("id", sid).
		Where("status", ShowStatStandard).
		Where("is_show", ShowOn).
		Find(&bInfo)

	if r.Error != nil {
		return nil
	}

	return bInfo
}


func GetActiveShows() []model.PrimaryKey {
	var f []*struct{
		ID model.PrimaryKey
	}
	Model().Select("id").Where("status = ?", ShowStatStandard).
		Order("id desc").
		Where("id IN ?", test.TestShowIds).
		Find(&f)

	activeIds := make([]model.PrimaryKey, len(f))
	for k, v := range f {
		activeIds[k] = v.ID
	}

	return activeIds
}

func GetActiveShowsName() []*ShowName {
	var f []*ShowName
	Model().Select("id", "name").Where("status = ?", ShowStatStandard).
		Order("id desc").
		Where("id IN ?", test.TestShowIds).
		Find(&f)

	return f
}

func GetActiveShowsWithType() []*Table {
	var f []*Table
	Model().Select("id", "show_type", "sub_show_type", "platform").
		Where("status = ?", ShowStatStandard).
		Where("id IN ?", test.TestShowIds).
		Order("id desc").
		Find(&f)

	return f
}


func GetOnShowing() []uint64 {
	var f []*Table
	Model().Select("id").
		Where("status = ?", ShowStatStandard).
		Where("release_at > 0 and (end_at = 0 or end_at >= ?)", time.Today() - 86400 * 14).
		Where("id IN ?", test.TestShowIds).
		Find(&f)

	sids := make([]uint64, len(f))

	for k, v := range f {
		sids[k] = v.ID
	}

	return sids
}

func GetShowInfo(sid uint64) Table {
	var f Table
	Model().Where("id", sid).
		Where("is_show", ShowOn).
		Limit(1).
		Find(&f)

	return f
}


func GetActiveShowsByType(sType int64, ssType int64) []Table {
	var f []Table
	mm := Model().Select("id", "show_type", "sub_show_type").
		Where("show_type", sType).
		Where("status = ?", ShowStatStandard).
		Where("is_show", ShowOn).
		Where("id IN ?", test.TestShowIds)

	if ssType >= 0 {
		mm = mm.Where("sub_show_type", ssType)
	}

	mm.Find(&f)
	return f
}


func GetShowsType(sid uint64) Table {
	var f Table
	Model().Select("id", "show_type", "sub_show_type").
		Where("id", sid).
		Find(&f)

	return f
}

// 存储上线时间
func StoreReleaseTime(rel uint, sid uint64){
	r := Model().Where("is_crawler_release = ?", ShowStatStandard).Limit(1).Find(nil, sid)
	if r.RowsAffected == 0 {
		return
	}

	r.Update("release_at", rel)
}