package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"iflytek.com/weipan4/learn-go/rag/configs"
)

var DB *gorm.DB

func InitDB() error {
	conf := configs.GetConfig()
	if conf == nil {
		return fmt.Errorf("config is not initialized")
	}
	mysqlConf := conf.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		mysqlConf.User,
		mysqlConf.Password,
		mysqlConf.Host,
		mysqlConf.Port,
		mysqlConf.DBName,
		mysqlConf.Charset,
		mysqlConf.ParseTime,
		mysqlConf.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}
	DB = db
	return nil
}
