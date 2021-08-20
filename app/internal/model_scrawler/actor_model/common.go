package actor_model

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

type ActorsName struct {
	Id model.PrimaryKey
	Name model.Varchar
	BirthYear model.Int
}
var actors []*ActorsName
func GetActor() []*ActorsName {

	var data []Table
	Model().Select("id", "name", "birth_year").
		Where("id IN ?", test.TestActorIds).
		Find(&data)

	actors = make([]*ActorsName, len(data))
	for  k, row := range data {
		actors[k] = &ActorsName{
			row.ID,
			row.Name,
			row.BirthYear,
		}
	}

	return actors
}


func GetActorName(aid uint64) string {
	// 有结果直接返回，防止对数据库的大量读取请求
	var data Table
	Model().Select("name").Where("id", aid).Find(&data)

	return data.Name
}


func GetNewActor() []*ActorsName {
	// 有结果直接返回，防止对数据库的大量读取请求
	if actors != nil {
		return actors
	}

	var data []Table
	Model().Select("id", "name", "birth_year").Find(&data)

	actors = make([]*ActorsName, len(data))
	for  k, row := range data {
		actors[k] = &ActorsName{
			row.ID,
			row.Name,
			row.BirthYear,
		}
	}

	return actors
}

// 保存演员
func SaveActor(name string) model.PrimaryKey {
	if name == ""{
		return 0
	}
	var row Table
	r := Model().Where("name = ?", name).Limit(1).Find(&row)
	if r.RowsAffected > 0 {
		return row.ID
	}

	d := &Table{Name: name}

	insert := Model().Create(&d)

	if insert.Error != nil {
		return 0
	}
	return d.ID
}