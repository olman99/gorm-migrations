package gormmigrations

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type command interface {
	execute() error
}

type migrate struct {
}

type migrateAdapter struct {
	migrate
	*gorm.DB
}

func newMigrate(db *gorm.DB) command {
	return &migrateAdapter{
		migrate: migrate{},
		DB:      db,
	}
}

func (c *migrate) execute(db *gorm.DB) error {
	for _, v := range migrationsFiles {
		if validateMigration(v.Name(), migrationsRun) {
			fmt.Println("Run migration " + v.Name())
			v.Run()
			migration := migration{Name: v.Name()}
			db.Create(&migration)
			fmt.Println("Migration " + v.Name() + " run successffully")
		}
	}

	return nil
}

func (c *migrateAdapter) execute() error {
	err := c.migrate.execute(c.DB)
	return err
}

type rollback struct {
}

type rollbackAdapter struct {
	rollback
	*gorm.DB
}

func newRollback(db *gorm.DB) command {
	return &rollbackAdapter{
		rollback: rollback{},
		DB:       db,
	}
}

func (c *rollback) execute(db *gorm.DB) error {
	for _, v := range migrationsFiles {
		if !validateMigration(v.Name(), migrationsRun) {
			fmt.Println("Drop migration " + v.Name())
			v.Drop()
			migration := migration{Name: v.Name()}
			db.Where("name = ?", v.Name()).Delete(&migration)
			fmt.Println("Migration " + v.Name() + " drop successffully")
		}
	}

	return nil
}

func (c *rollbackAdapter) execute() error {
	err := c.rollback.execute(c.DB)
	return err
}

type rollbackLast struct {
}

type rollbackLastAdapter struct {
	rollbackLast
	*gorm.DB
}

func newRollbackLast(db *gorm.DB) command {
	return &rollbackLastAdapter{
		rollbackLast: rollbackLast{},
		DB:           db,
	}
}

func (c *rollbackLast) execute(db *gorm.DB) error {
	index := len(migrationsRun) - 1
	indexAll := len(migrationsFiles) - 1

	for {
		if migrationsFiles[indexAll].Name() == migrationsRun[index].Name {
			fmt.Println("Drop migration " + migrationsFiles[indexAll].Name())
			migrationsFiles[indexAll].Drop()
			migration := migration{Name: migrationsFiles[indexAll].Name()}
			db.Where("name = ?", migrationsFiles[indexAll].Name()).Delete(&migration)
			fmt.Println("Migration " + migrationsFiles[indexAll].Name() + " drop successffully")
			break
		}
		indexAll--
	}

	return nil
}

func (c *rollbackLastAdapter) execute() error {
	err := c.rollbackLast.execute(c.DB)
	return err
}

type makeMigrations struct {
	name  string
	table string
	path  string
}

func newMakeMigrations(name string, table string, path string) command {
	return &makeMigrations{
		name:  name,
		table: table,
		path:  path,
	}
}

func (c *makeMigrations) execute() error {

	t := time.Now()

	name := fmt.Sprintf("m%d%02d%02d%02d%02d_%s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), c.name)

	contentString := createMigrationTemplate(name, c.table)

	err := writeFile(c.path+name, contentString)

	if err != nil {
		return err
	}

	return nil
}

func getCommand(args *cmdArgs, db *gorm.DB) command {
	switch args.command {
	case "migrations:migrate":
		return newMigrate(db)
	case "migrations:rollback":
		return newRollback(db)
	case "migrations:rollback:last":
		return newRollbackLast(db)
	case "migrations:create":
		return newMakeMigrations(args.params["migration"], args.flags["table"], args.flags["path"])
	case "migrations:seeder":
	}

	return nil
}
