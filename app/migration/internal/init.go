package internal

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"guduo/app/internal/config"
	"guduo/pkg/db"
)

func InitDB() {
	// 初始化mysql
	mysqlErr := db.InitMysqlConn(config.MysqlCrawlerConfig, "crawler")
	if mysqlErr != nil {
		logrus.Error(fmt.Sprintf("mysql连接初始化失败：%s", mysqlErr))
	}

	mysqlErr = db.InitMysqlConn(config.MysqlCleanConfig, "clean")
	if mysqlErr != nil {
		logrus.Error(fmt.Sprintf("mysql连接初始化失败：%s", mysqlErr))
	}

	mysqlErr = db.InitMysqlConn(config.MysqlLolipopConfig, "lolipop")
	if mysqlErr != nil {
		logrus.Error(fmt.Sprintf("mysql连接初始化失败：%s", mysqlErr))
	}

	mysqlErr = db.InitMysqlConn(config.MysqlCarlConfig, "carl")
	if mysqlErr != nil {
		logrus.Error(fmt.Sprintf("mysql连接初始化失败：%s", mysqlErr))
	}
}
