package orm

import (
	"fmt"
	"log"
	"strings"
)

type DatabaseType = string

const (
	DatabaseTypeSQLite   DatabaseType = "sqlite"
	DatabaseTypePostgres DatabaseType = "postgres"
	DatabaseTypeMysql    DatabaseType = "mysql"

	MysqlDefaultPort    = 3306
	PostgresDefaultPort = 5432
)

type Config struct {
	Type     DatabaseType `mapstructure:"type"`
	Host     string       `mapstructure:"host"`
	Port     int          `mapstructure:"port"`
	User     string       `mapstructure:"user"`
	Password string       `mapstructure:"password"`
	Name     string       `mapstructure:"name"`
	File     string       `mapstructure:"file"`
	SSLMode  string       `mapstructure:"ssl_mode"`
	DSN      string       `mapstructure:"dsn"`
}

// GetDSN get connection string	and type
func (d Config) GetDSN() (t DatabaseType, dsn string) {
	t = d.Type
	dsn = d.DSN
	if dsn == "" {
		switch t {
		case DatabaseTypeSQLite:
			if !(strings.HasSuffix(d.File, ".db") && len(d.File) > 3) {
				log.Fatalf("db name error.")
			}
			dsn = fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental", d.File)
		case DatabaseTypePostgres:
			if d.Password != "" {
				dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Taipei",
					d.Host, d.User, d.Password, d.Name, d.Port, d.SSLMode)
			} else {
				dsn = fmt.Sprintf("host=%s user=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Taipei",
					d.Host, d.User, d.Name, d.Port, d.SSLMode)
			}
		case DatabaseTypeMysql:
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
				d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
		}
	}
	return
}

func GetDefaultConfig(t DatabaseType) Config {
	switch t {
	case DatabaseTypePostgres:
		return Config{
			Type:     DatabaseTypePostgres,
			Host:     "localhost",
			Port:     PostgresDefaultPort,
			User:     "user",
			Password: "password",
			Name:     "database",
			SSLMode:  "disable",
		}
	case DatabaseTypeMysql:
		return Config{
			Type:     DatabaseTypeMysql,
			Host:     "localhost",
			Port:     MysqlDefaultPort,
			User:     "user",
			Password: "password",
			Name:     "database",
			SSLMode:  "disable",
		}
	default:
		return Config{
			Type: DatabaseTypeSQLite,
			File: "database.db",
		}
	}
}
