package show_detail_model

import (
	"gorm.io/gorm"
	"guduo/app/internal/test"
	"guduo/pkg/db"
	"guduo/pkg/model"
)

var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

type DetailUrls []*struct {
	ShowId model.PrimaryKey
	Url    model.Text
}

type DetailUrl struct {
	ShowId model.PrimaryKey
	Url    model.Text
}

func GetDetailUrl(platformId model.ForeignKey, showIds []model.ForeignKey) DetailUrls {
	if len(showIds) == 0 {
		return nil
	}
	var d DetailUrls
	Model().Select("show_id", "IF(true_url is null or true_url = '', url, true_url) as url").
		Where("platform_id = ?", platformId).
		Where("show_id IN ?", showIds).
		Where("usable", 1).
		Order("show_id desc").
		Find(&d)

	return d
}

func GetDetailUrlNew(platformId model.ForeignKey) DetailUrls {
	tableName := "`show`"
	originSql := `select t1.show_id,
						   IF(t1.true_url is null or t1.true_url = '', t1.url, t1.true_url) as url
					from show_detail t1,
						 ` + tableName + ` t2
					where t1.platform_id = ?
					  and t1.usable = 1
					  and t1.show_id = t2.id
					  and t2.status = 1
					  and t1.show_id IN ?
						order by t1.show_id desc`

	var d DetailUrls
	Model().Raw(originSql, platformId, test.TestShowIds).Scan(&d)

	//Model().Select("show_id", "IF(true_url is null or true_url = '', url, true_url) as url").
	//	Where("platform_id = ?", platformId).
	//	Where("show_id IN ?", showIds).
	//	Where("usable", 1).
	//	Order("show_id desc").
	//	Find(&d)

	return d
}

func SaveTrueUrl(url string, sid uint64, pid uint64) {
	if url != "" {
		Model().Where("show_id", sid).
			Where("platform_id", pid).
			Update("true_url", url)
	} else {
		Model().Where("show_id", sid).
			Where("platform_id", pid).
			Update("usable", 0)
	}
}

func ReportErrorUrl(u string) {
	Model().Where("url = ? or true_url = ?", u, u).
		Update("usable", 0)
}
