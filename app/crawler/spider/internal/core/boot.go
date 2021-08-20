package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"guduo/app/crawler/spider/internal/core/proxy"
	config2 "guduo/app/internal/config"
	"guduo/pkg/db"
	"guduo/pkg/log"
)

func Init() {
	// 初始化日志系统
	log.InitLogger()

	// 初始化mysql和redis
	//_ = db.InitRedisConn(config.RedisConfig)
	mysqlErr := db.InitMysqlConn(config2.MysqlCrawlerConfig, "crawler")
	if mysqlErr != nil {
		logrus.Error(fmt.Sprintf("mysql连接初始化失败：%s", mysqlErr))
		fmt.Println(mysqlErr)
		return
	}

	mysqlErr = db.InitMysqlConn(config2.MysqlCleanConfig, "clean")
	if mysqlErr != nil {
		logrus.Error(fmt.Sprintf("mysql连接初始化失败：%s", mysqlErr))
		fmt.Println(mysqlErr)
		return
	}

	// 初始化代理池
	proxy.InitProxyPool()

	// 初始化app
	//initApp()
}
