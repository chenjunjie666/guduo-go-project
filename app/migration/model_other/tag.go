package model_other

import (
	db2 "guduo/pkg/db"
	"guduo/pkg/model"
)

type SST struct {
	ShowId uint64
	ThemeId uint64
}

func (s SST) TableName() string {
	return "show_select_theme"
}

type ST struct {
	ID model.PrimaryKey
	Name string
}

func (s ST) TableName() string {
	return "show_theme"
}

var sst = make(map[uint64][]string)

func GetShowTheme(sid uint64) []string {
	if len(sst) > 0 {
		if _, ok := sst[sid]; !ok{
			return []string{}
		}
		return sst[sid]
	}else{
		var dids []SST
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SST{}).Select("theme_id", "show_id").Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.ThemeId)
		}

		var dn []ST
		db.Model(&ST{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			sst[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						sst[sid2] = append(sst[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := sst[sid]; !ok{
			return []string{}
		}
		return sst[sid]
	}
}



type SSTN struct {
	ShowId uint64
	ThemeId uint64
}

func (s SSTN) TableName() string {
	return "show_select_theme_new"
}

type STN struct {
	ID model.PrimaryKey
	Name string
}

func (s STN) TableName() string {
	return "show_theme_new"
}

var sstn = make(map[uint64][]string)

func GetShowThemeNew(sid uint64) []string {
	if len(sstn) > 0 {
		if _, ok := sstn[sid]; !ok{
			return []string{}
		}
		return sstn[sid]
	}else{
		var dids []SSTN
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SSTN{}).Select("theme_id", "show_id").Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.ThemeId)
		}

		var dn []STN
		db.Model(&STN{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			sstn[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						sstn[sid2] = append(sstn[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := sstn[sid]; !ok{
			return []string{}
		}
		return sstn[sid]
	}
}

