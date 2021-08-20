package model_other

import (
	db2 "guduo/pkg/db"
	"guduo/pkg/model"
)

type SS struct {
	ShowId uint64
	ScriptwriterId uint64
}

func (s SS) TableName() string {
	return "show_scriptwriter"
}

type S struct {
	ID model.PrimaryKey
	Name string
}

func (s S) TableName() string {
	return "scriptwriter"
}

var ss = make(map[uint64][]string)

func GetScriptwriter(sid uint64) []string {
	if len(ss) > 0 {
		if _, ok := ss[sid]; !ok{
			return []string{}
		}
		return ss[sid]
	}else{
		var dids []SS
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SS{}).Select("scriptwriter_id", "show_id").Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.ScriptwriterId)
		}

		var dn []S
		db.Model(&S{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			ss[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						ss[sid2] = append(ss[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := ss[sid]; !ok{
			return []string{}
		}
		return ss[sid]
	}
}











type SP struct {
	ShowId uint64
	ProducerId uint64
}

func (s SP) TableName() string {
	return "show_producer"
}

type P struct {
	ID model.PrimaryKey
	Name string
}

func (s P) TableName() string {
	return "producer"
}

var p = make(map[uint64][]string)

func GetProducer(sid uint64) []string {
	if len(p) > 0 {
		if _, ok := p[sid]; !ok{
			return []string{}
		}
		return p[sid]
	}else{
		var dids []SP
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SP{}).Select("producer_id", "show_id").Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.ProducerId)
		}

		var dn []P
		db.Model(&P{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			p[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						p[sid2] = append(p[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := p[sid]; !ok{
			return []string{}
		}
		return p[sid]
	}
}





type SPUB struct {
	ShowId uint64
	PublisherId uint64
}

func (s SPUB) TableName() string {
	return "show_publisher"
}

type PUB struct {
	ID model.PrimaryKey
	Name string
}

func (s PUB) TableName() string {
	return "publisher"
}

var sp = make(map[uint64][]string)
func GetPublisher(sid uint64) []string {
	if len(sp) > 0 {
		if _, ok := sp[sid]; !ok{
			return []string{}
		}
		return sp[sid]
	}else{
		var dids []SPUB
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SPUB{}).Select("publisher_id", "show_id").Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.PublisherId)
		}

		var dn []PUB
		db.Model(&PUB{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			sp[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						sp[sid2] = append(sp[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := sp[sid]; !ok{
			return []string{}
		}
		return sp[sid]
	}
}




type SC struct {
	ShowId uint64
	CompanyId uint64
}

func (s SC) TableName() string {
	return "company_show_relation"
}

type C struct {
	ID model.PrimaryKey
	Name string
}

func (s C) TableName() string {
	return "company"
}

var c1 = make(map[uint64][]string)
func GetCompany1(sid uint64) []string {
	if len(c1) > 0 {
		if _, ok := c1[sid]; !ok{
			return []string{}
		}
		return c1[sid]
	}else{
		var dids []SC
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SC{}).Select("company_id", "show_id").Where("relation_type", 1).
			Where("status != -1").
			Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.CompanyId)
		}

		var dn []C
		db.Model(&C{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			c1[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						c1[sid2] = append(c1[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := c1[sid]; !ok{
			return []string{}
		}
		return c1[sid]
	}
}


var c2 = make(map[uint64][]string)
func GetCompany2(sid uint64) []string {
	if len(c2) > 0 {
		if _, ok := c2[sid]; !ok{
			return []string{}
		}
		return c2[sid]
	}else{
		var dids []SC
		db := db2.GetLoliPopMysqlConn()
		db.Model(&SC{}).Select("company_id", "show_id").Where("relation_type", 2).
			Where("status != -1").
			Find(&dids)
		if len(dids) == 0 {
			return []string{}
		}

		tmp := make(map[uint64][]uint64)
		for _, v := range dids {
			if _, ok := tmp[v.ShowId]; !ok {
				tmp[v.ShowId] = make([]uint64, 0, 10)
			}
			tmp[v.ShowId] = append(tmp[v.ShowId], v.CompanyId)
		}

		var dn []C
		db.Model(&C{}).Select("id", "name").Find(&dn)

		for sid2, ssids := range tmp{
			c2[sid2] = make([]string, 0, 10)
			for _, ssid := range ssids {
				for _, row := range dn {
					if row.ID == ssid {
						c2[sid2] = append(c2[sid2], row.Name)
					}
				}
			}
		}

		if _, ok := c2[sid]; !ok{
			return []string{}
		}
		return c2[sid]
	}
}
