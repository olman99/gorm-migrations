package gormmigrations

import "fmt"

func createMigrationTemplate(name string, table string) string {
	return fmt.Sprintf(
		`package migrations

import (
	gormmigrations "github.com/olman99/gorm-migrations"
	"gorm.io/gorm"
)
		
type %s struct {
	DB *gorm.DB
}

type %s_Table struct {

}

func (%s_Table) TableName() string {
	return "%s"
}
		
func New%s(db *gorm.DB) gormmigrations.IMigration {
	return &%s{DB:db,}
}
		
func (m *%s) Run(){}
		
func (m *%s) Drop(){}
		
func (m *%s) Name() string {
	return "%s"
}
	`, name, name, name, table, name, name, name, name, name, name,
	)
}
