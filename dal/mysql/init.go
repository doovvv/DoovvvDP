package mysql

import (
	"doovvvDP/config"
	"fmt"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
var (
	DB  *gorm.DB
	err error
)
func Init() {
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
		},
	)
	if err != nil {
		panic(err)
	}
}