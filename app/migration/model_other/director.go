package model_other

import (
	db2 "guduo/pkg/db"
	"guduo/pkg/model"
)

type ShowDir struct {
	ShowId uint64
	DirectorId uint64
}

func (d ShowDir) TableName() string {
	return "show_director"
}


type Dir struct {
	ID model.PrimaryKey
	Name model.Varchar
}
func (d Dir) TableName() string {
	return "director"
}

var dir = make(map[uint64][]string)
func GetDirector(sid uint64) []string {
	if len(dir) > 0 {
		if _, ok := dir[sid]; !ok{
			return []string{}
		}
		return dir[sid]
	}else{
		var dids []ShowDir
		db := db2.GetLoliPopMysqlConn()
		db.Model(&ShowDir{}).Select("director_id", "show_id").Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.DirectorId)
		}

		var dn []Dir
		db.Model(&Dir{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			dir[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						dir[sid2] = append(dir[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := dir[sid]; !ok{
			return []string{}
		}
		return dir[sid]
	}
}


