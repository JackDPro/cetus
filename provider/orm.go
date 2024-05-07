package provider

import (
	"cetus/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type Orm struct {
	Db *gorm.DB
}

var ormInstance *Orm
var ormOnce sync.Once

func GetOrm() *Orm {
	ormOnce.Do(func() {
		dbConf := config.GetDatabaseConfig()
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		ormInstance = &Orm{
			Db: db,
		}
	})
	return ormInstance
}
