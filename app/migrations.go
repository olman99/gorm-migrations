package gormmigrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type migration struct {
	Name string
	Date *time.Time
}

func (migration) TableName() string {
	return "migrations"
}

type IMigration interface {
	Name() string
	Run()
	Drop()
}

var migrationsFiles []IMigration
var migrationsRun []migration

func getMigrations(db *gorm.DB) []migration {

	if db.Migrator().HasTable("migrations") {
		err := db.Find(&migrationsRun).Error
		if err != nil {
			fmt.Println("Can't loaded table migrations")
			panic(err)
		}

	} else {
		fmt.Println("Create table migrations")
		db.Migrator().CreateTable(&migration{})
	}

	return migrationsRun
}

func validateMigration(name string, migrations []migration) bool {
	for _, v := range migrations {
		if v.Name == name {
			return false
		}
	}
	return true
}
