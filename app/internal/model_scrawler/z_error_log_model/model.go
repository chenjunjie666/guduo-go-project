package z_error_log_model

import (
	"gorm.io/gorm"
	"guduo/pkg/db"
	"guduo/pkg/model"
)

type Table struct {
	model.Fields
	ShowId model.ForeignKey
	ActorId model.ForeignKey
	PlatformId model.ForeignKey
	JobName model.Varchar
	JobDesc model.Varchar
	Int1 model.Int
	Int2 model.Int
	F1 model.Float
	F2 model.Float
	JobAt model.Int
	Desc model.Varchar
}

func (d Table) TableName() string {
	return "z_error_log"
}

// 直接获取当前模型的对象
var m *gorm.DB

func Model() *gorm.DB {
	if m == nil {
		m = db.GetCrawlerMysqlConn()
	}
	return m.Model(&Table{})
}

func ShowDBError (sid uint64, err error, pid uint64, job string, jobDesc string) {
	errMsg := err.Error()

	res := &Table{
		ShowId: sid,
		PlatformId: pid,
		JobName: job,
		JobDesc: jobDesc,
		Desc: errMsg,
	}
	Model().Create(res)
}
