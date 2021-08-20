package carl_shows_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/model_scrawler/show_detail_model"
	"guduo/pkg/db"
	"time"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCarlMysqlConn()
	}
	return m.Model(&Table{})
}

func Sync () {
	var tmpRes *show_detail_model.Table
	show_detail_model.Model().Select("show_id").Order("show_id desc").Limit(1).Find(&tmpRes)
	maxId := tmpRes.ID

	show_detail_model.Model().Where("show_id > ?", maxId).Delete(nil)
	time.Sleep(time.Second * 10)
	tn := show_detail_model.Table{}.TableName()
	db.GetCrawlerMysqlConn().Exec("alter table `" + tn + "` AUTO_INCREMENT=1;")

	var res []*Table
	Model().Select("platform_id", "linked_id", "url").Where("depth", 1).
		Where("linked_id > ?", maxId).
		Find(&res)

	arr := make([]*show_detail_model.Table, 0, 400)
	for _, v := range res {
		row := &show_detail_model.Table{
			ShowId:     v.LinkedId,
			PlatformId: v.PlatformId,
			Url:        v.Url,
			Usable:     1,
			TrueUrl:    "",
		}
		arr = append(arr, row)
		if len(arr) >= 400 {
			show_detail_model.Model().Create(&arr)
			arr = arr[:0]
		}
	}

	if len(arr) > 0 {
		show_detail_model.Model().Create(&arr)
	}
}
