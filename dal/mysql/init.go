package mysql

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"doovvvDP/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	logFile, _ := os.OpenFile("logs/mysql.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags), // io writer（日志输出的目标，如：文件、os.Stdout等）
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 日志级别
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MyConfig.MySQLConfig.User,
		config.MyConfig.MySQLConfig.Password,
		config.MyConfig.MySQLConfig.Host,
		strconv.Itoa(config.MyConfig.MySQLConfig.Port),
		config.MyConfig.MySQLConfig.DatabaseName,
	)
	fmt.Println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger:                 newLogger,
		},
	)
	if err != nil {
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(config.MyConfig.MySQLConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MyConfig.MySQLConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MyConfig.MySQLConfig.MaxIdleTime) * time.Second)
}
