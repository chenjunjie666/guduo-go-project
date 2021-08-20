package model_other

import (
	db2 "guduo/pkg/db"
)

type ShowPlt struct {
	ShowId uint64
	PlatformId uint64
}

func (p ShowPlt) TableName() string {
	return "show_platform"
}

var plat = make(map[uint64][]uint64)

func GetPlat(sid uint64) []uint64 {
	if len(plat) > 0 {
		if _, ok := plat[sid]; !ok{
			return []uint64{}
		}
		return plat[sid]
	}else{
		var dids []ShowPlt
		db := db2.GetLoliPopMysqlConn()
		db.Model(&ShowPlt{}).Select("platform_id", "show_id").Find(&dids)

		if len(dids) == 0 {
			return []uint64{}
		}

		for _, v := range dids {
			if _, ok := plat[v.ShowId]; !ok {
				plat[v.ShowId] = make([]uint64, 0, 5)
			}
			plat[v.ShowId] = append(plat[v.ShowId], v.PlatformId)
		}
		if _, ok := plat[sid]; !ok{
			return []uint64{}
		}
		return plat[sid]
	}
}