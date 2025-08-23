package orm_test

import (
	"context"
	"testing"
	"time"

	orm "github.com/Wuchieh/go-server-orm"
	"gorm.io/gorm"
)

func TestNew(t *testing.T) {
	db, err := orm.New(orm.Config{
		Type: orm.DatabaseTypeSQLite,
		DSN:  "file::memory:?cache=shared",
	})

	if err != nil {
		t.Fatalf("New err: %v", err)
	}

	type user struct {
		gorm.Model
		Name string
	}

	err = db.AutoMigrate(&user{})
	if err != nil {
		t.Fatalf("AutoMigrate err: %v", err)
	}

	tempName := time.Now().String()
	db.Create(&user{
		Name: tempName,
	})

	var record user
	db.First(&record)

	if record.Name != tempName {
		t.Fatalf("name not match")
	}
}

func TestSetup(t *testing.T) {
	err := orm.Setup(orm.Config{
		Type: orm.DatabaseTypeSQLite,
		DSN:  "file::memory:?cache=shared",
	})

	if err != nil {
		t.Fatalf("New err: %v", err)
	}

	defer func() {
		err := orm.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	db := orm.GetDB()

	type user struct {
		gorm.Model
		Name string
	}

	err = db.AutoMigrate(&user{})
	if err != nil {
		t.Fatalf("AutoMigrate err: %v", err)
	}

	tempName := time.Now().String()
	db.Create(&user{
		Name: tempName,
	})

	var record user
	db.First(&record)

	if record.Name != tempName {
		t.Fatalf("name not match")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	db = orm.GetDBWithContent(ctx)
	if err := db.
		Create(&user{Name: time.Now().String()}).
		Error; err == nil {
		t.Fatalf("context not work")
	}
}
