package config

import (
	"guduo/app/internal/boot"
	"guduo/pkg/db"
)

var RedisConfig = &db.RedisConfig{
	Addr: boot.Cfg.Redis.Host,
	Pass: boot.Cfg.Redis.Pass, // no password set
	DB:   boot.Cfg.Redis.DB,   // use default DB
	Port: boot.Cfg.Redis.Port,
}
