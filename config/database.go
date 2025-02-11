package config

import (
	"os"
	"sync"
)

type DatabaseConf struct {
	Host            string
	Port            string
	Database        string
	DeviceDatabase  string
	Username        string
	Password        string
	MigrateSelfOnly bool
}

var databaseConfigInstance *DatabaseConf
var databaseConfigOnce sync.Once

func GetDatabaseConfig() *DatabaseConf {
	databaseConfigOnce.Do(func() {
		host := "127.0.0.1"
		if os.Getenv("DB_HOST") != "" {
			host = os.Getenv("DB_HOST")
		}
		port := "3306"
		if os.Getenv("DB_PORT") != "" {
			port = os.Getenv("DB_PORT")
		}
		database := "test"
		if os.Getenv("DB_DATABASE") != "" {
			database = os.Getenv("DB_DATABASE")
		}
		username := "root"
		if os.Getenv("DB_USERNAME") != "" {
			username = os.Getenv("DB_USERNAME")
		}
		password := "password"
		if os.Getenv("DB_PASSWORD") != "" {
			password = os.Getenv("DB_PASSWORD")
		}
		migrateSelfOnly := false
		if os.Getenv("DB_MIGRATE_SELF_ONLY") == "true" {
			migrateSelfOnly = true
		}
		databaseConfigInstance = &DatabaseConf{
			Host:            host,
			Port:            port,
			Database:        database,
			Username:        username,
			Password:        password,
			MigrateSelfOnly: migrateSelfOnly,
		}
	})
	return databaseConfigInstance
}
