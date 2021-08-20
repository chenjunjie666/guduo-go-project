package db

import (
	"fmt"
	"gorm.io/gorm/logger"
	"guduo/pkg/errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormDB struct {
	Crawler *gorm.DB
	Clean   *gorm.DB
	LoliPop *gorm.DB
	Carl    *gorm.DB
}

func InitMysqlConn(opt *MysqlConfig, name string) error {
	if gormDB.Crawler != nil && gormDB.Clean != nil && gormDB.LoliPop != nil && gormDB.Carl != nil {
		return nil
	}
	_gormDB, err := NewMysqlConn(opt)
	if name == "crawler"{
		gormDB.Crawler = _gormDB
	}else if name == "clean"{
		gormDB.Clean = _gormDB
	}else if name == "lolipop"{
		gormDB.LoliPop = _gormDB
	}else if name == "carl"{
		gormDB.Carl = _gormDB
	}else {
		return errors.AppError("mysql", "mysql配置不正确")
	}
	return err
}

func NewMysqlConn(opt *MysqlConfig) (*gorm.DB, error) {
	dns, err := opt.toDNS()
	if err != nil {
		return nil, err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second * 5,   // 慢 SQL 阈值
			LogLevel:      logger.Warn, // Log level
		},
	)

	cfg := &gorm.Config{
		Logger: newLogger,
	}

	db, err := gorm.Open(mysql.Open(dns), cfg)

	if err != nil {
		return nil, err
	}

	mysqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	mysqlDB.SetMaxIdleConns(20)               // 最大空闲的链接
	mysqlDB.SetMaxOpenConns(150)              // 链接池最大链接数
	mysqlDB.SetConnMaxLifetime(time.Second * 120) // 单个链接的最大链接时长为2分钟

	return db, nil
}

func GetCrawlerMysqlConn() *gorm.DB {
	return gormDB.Crawler
}

func GetCleanMysqlConn() *gorm.DB {
	return gormDB.Clean
}

func GetLoliPopMysqlConn() *gorm.DB {
	return gormDB.LoliPop
}

func GetCarlMysqlConn() *gorm.DB {
	return gormDB.Carl
}

type MysqlConfig struct {
	Host    string
	User    string
	Pass    string
	DbName  string
	Port    int
	Charset string
}

func (m *MysqlConfig) init() {
	if m.Port == 0 {
		m.Port = 3306
	}

	if m.Charset == "" {
		m.Charset = "utf8mb4"
	}
}

func (m *MysqlConfig) toDNS() (string, error) {
	m.init()

	c := m.check()
	if c != "" {
		return "", errors.AppError("mysql", fmt.Sprintf("缺少链接参数：%s", c))
	}

	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		m.User,
		m.Pass,
		m.Host,
		m.Port,
		m.DbName,
		m.Charset,
	)
	return dns, nil
}

func (m *MysqlConfig) check() string {
	if m.User == "" {
		return "User不能为空"
	}
	if m.Host == "" {
		return "Host不能为空"
	}
	if m.DbName == "" {
		return "DbName不能为空"
	}

	return ""
}
