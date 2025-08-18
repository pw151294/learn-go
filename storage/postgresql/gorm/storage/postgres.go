package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"iflytek.com/weipan4/learn-go/storage/postgresql/gorm/config"
)

const postgresDsn = "user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai"

var sqlDB *gorm.DB

func GetDB() *gorm.DB {
	return sqlDB
}

func InitStorage() error {
	dbCfg := config.GetDataBaseConfig()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf(postgresDsn, dbCfg.User, dbCfg.Password, dbCfg.Database, dbCfg.Port, dbCfg.Sslmode),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}))
	if err != nil {
		return err
	}
	sqlDB = db
	return nil
}
