package provider

import (
	"fmt"
	"github.com/JackDPro/cetus/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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

		// Build GORM configuration
		conf := &gorm.Config{}
		if config.GetAppConfig().Env == "dev" {
			conf.Logger = logger.Default.LogMode(logger.Info)
		}
		conf.DisableForeignKeyConstraintWhenMigrating = dbConf.MigrateSelfOnly
		conf.IgnoreRelationshipsWhenMigrating = dbConf.MigrateSelfOnly

		// Select database driver based on type
		var dialector gorm.Dialector
		switch dbConf.Type {
		case "postgres":
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				dbConf.Host, dbConf.Username, dbConf.Password, dbConf.Database, dbConf.Port, dbConf.SSLMode)
			dialector = postgres.Open(dsn)
		case "mysql":
			fallthrough
		default:
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
			dialector = mysql.Open(dsn)
		}

		db, err := gorm.Open(dialector, conf)
		if err != nil {
			panic(fmt.Sprintf("failed to connect database: %v", err))
		}
		ormInstance = &Orm{
			Db: db,
		}
	})
	return ormInstance
}
