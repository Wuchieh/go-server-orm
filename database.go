package orm

import (
	"context"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func New(cfg Config, opts ...gorm.Option) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	t, dsn := cfg.GetDSN()

	switch t {
	case DatabaseTypeSQLite:
		db, err = gorm.Open(sqlite.Open(dsn), opts...)
	case DatabaseTypePostgres:
		db, err = gorm.Open(postgres.Open(dsn), opts...)
	case DatabaseTypeMysql:
		db, err = gorm.Open(mysql.Open(dsn), opts...)
	}

	return db, err
}

func Setup(cfg Config, opts ...gorm.Option) error {
	t, dsn := cfg.GetDSN()

	var err error
	switch t {
	case DatabaseTypeSQLite:
		db, err = gorm.Open(sqlite.Open(dsn), opts...)
	case DatabaseTypePostgres:
		db, err = gorm.Open(postgres.Open(dsn), opts...)
	case DatabaseTypeMysql:
		db, err = gorm.Open(mysql.Open(dsn), opts...)
	}

	return err
}

func GetDB(tx ...*gorm.DB) *gorm.DB {
	db := db

	for _, t := range tx {
		if t != nil {
			db = t
		}
	}

	return db
}

func GetDBWithContent(ctx context.Context, tx ...*gorm.DB) *gorm.DB {
	db := GetDB(tx...)

	if ctx == nil {
		ctx = context.Background()
	}

	return db.WithContext(ctx)
}

func Close() error {
	db := GetDB()
	if db == nil {
		return nil
	}

	s, err := db.DB()
	if err != nil {
		return err
	}

	return s.Close()
}
