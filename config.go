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
)

type Config struct {
	Type        DatabaseType `mapstructure:"type"`
	Host        string       `mapstructure:"host"`
	Port        int          `mapstructure:"port"`
	User        string       `mapstructure:"user"`
	Password    string       `mapstructure:"password"`
	Name        string       `mapstructure:"name"`
	File        string       `mapstructure:"file"`
	TablePrefix string       `mapstructure:"table_prefix"`
	SSLMode     string       `mapstructure:"ssl_mode"`
	DSN         string       `mapstructure:"dsn"`
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
