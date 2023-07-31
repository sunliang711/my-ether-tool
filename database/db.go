package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Conn *gorm.DB
)

func InitDB(level string, dsn string) {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,        // 慢 SQL 阈值
			LogLevel:                  ormLogLevel(level), // 日志级别
			IgnoreRecordNotFoundError: true,               // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,               // 彩色打印
		},
	)
	Conn, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	Conn.AutoMigrate(&Account{})
	Conn.AutoMigrate(&Network{})
}

func ormLogLevel(levelString string) logger.LogLevel {
	switch levelString {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}
	return logger.Silent
}
