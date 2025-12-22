package config

import (
	"os"
	"sync"
)

type DatabaseConf struct {
	Type            string // Database type: mysql, postgres
	Host            string
	Port            string
	Database        string
	DeviceDatabase  string
	Username        string
	Password        string
	SSLMode         string // PostgreSQL SSL mode: disable, require, verify-ca, verify-full
	MigrateSelfOnly bool
}

var databaseConfigInstance *DatabaseConf
var databaseConfigOnce sync.Once

func GetDatabaseConfig() *DatabaseConf {
	databaseConfigOnce.Do(func() {
		dbType := "mysql"
		if os.Getenv("DB_TYPE") != "" {
			dbType = os.Getenv("DB_TYPE")
		}
		host := "127.0.0.1"
		if os.Getenv("DB_HOST") != "" {
			host = os.Getenv("DB_HOST")
		}
		// Default port based on database type
		defaultPort := "3306"
		if dbType == "postgres" {
			defaultPort = "5432"
		}
		port := defaultPort
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
		sslMode := "disable"
		if os.Getenv("DB_SSLMODE") != "" {
			sslMode = os.Getenv("DB_SSLMODE")
		}
		migrateSelfOnly := false
		if os.Getenv("DB_MIGRATE_SELF_ONLY") == "true" {
			migrateSelfOnly = true
		}
		databaseConfigInstance = &DatabaseConf{
			Type:            dbType,
			Host:            host,
			Port:            port,
			Database:        database,
			Username:        username,
			Password:        password,
			SSLMode:         sslMode,
			MigrateSelfOnly: migrateSelfOnly,
		}
	})
	return databaseConfigInstance
}
