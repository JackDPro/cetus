package provider

import (
	"fmt"
	"github.com/JackDPro/cetus/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		conf := &gorm.Config{}
		if config.GetAppConfig().Debug {
			conf.Logger = logger.Default.LogMode(logger.Info)
		}
		conf.DisableForeignKeyConstraintWhenMigrating = dbConf.MigrateSelfOnly
		db, err := gorm.Open(mysql.Open(dsn), conf)
		if err != nil {
			panic("failed to connect database")
		}
		ormInstance = &Orm{
			Db: db,
		}
	})
	return ormInstance
}
