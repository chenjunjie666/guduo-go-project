package config

import (
	boot2 "guduo/app/internal/boot"
	"guduo/pkg/db"
)

var MysqlCrawlerConfig = &db.MysqlConfig{
	Host:   boot2.Cfg.Database.Crawler.Host,
	User:   boot2.Cfg.Database.Crawler.User,
	Pass:   boot2.Cfg.Database.Crawler.Pass,
	DbName: boot2.Cfg.Database.Crawler.Dbname,
	Port:   boot2.Cfg.Database.Crawler.Port,
}


var MysqlCleanConfig = &db.MysqlConfig{
	Host:   boot2.Cfg.Database.Clean.Host,
	User:   boot2.Cfg.Database.Clean.User,
	Pass:   boot2.Cfg.Database.Clean.Pass,
	DbName: boot2.Cfg.Database.Clean.Dbname,
	Port:   boot2.Cfg.Database.Clean.Port,
}



var MysqlLolipopConfig = &db.MysqlConfig{
	Host:   boot2.Cfg.Database.Lolipop.Host,
	User:   boot2.Cfg.Database.Lolipop.User,
	Pass:   boot2.Cfg.Database.Lolipop.Pass,
	DbName: boot2.Cfg.Database.Lolipop.Dbname,
	Port:   boot2.Cfg.Database.Lolipop.Port,
}


var MysqlCarlConfig = &db.MysqlConfig{
	Host:   boot2.Cfg.Database.Carl.Host,
	User:   boot2.Cfg.Database.Carl.User,
	Pass:   boot2.Cfg.Database.Carl.Pass,
	DbName: boot2.Cfg.Database.Carl.Dbname,
	Port:   boot2.Cfg.Database.Carl.Port,
}