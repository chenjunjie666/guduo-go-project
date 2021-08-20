package boot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"guduo/app/internal/config"
	"guduo/app/weixin_mini/internal/core"
	"guduo/app/weixin_mini/internal/router"
	"guduo/pkg/db"
	"guduo/pkg/log"
)

func Init() {
	// 初始化日志系统
	log.InitLogger()
	InitDB()
	//InitRedis()
	s := InitServer()
	RunServer(s)
}

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
}

func InitRedis() {
	// 初始化redis
	RedisErr := db.InitRedisConn(config.RedisConfig)
	if RedisErr != nil {
		logrus.Error(fmt.Sprintf("redis连接初始化失败：%s", RedisErr))
		fmt.Println(RedisErr)
	}
}

func InitServer() *gin.Engine {
	s := core.GetServer()

	router.InitRouter()

	return s
}

func RunServer(s *gin.Engine) {
	e := s.Run("0.0.0.0:90")

	if e != nil {
		fmt.Println(e)
		panic(fmt.Sprintf("启动http服务出错:%s", e))
	}
}
