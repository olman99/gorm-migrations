package gormmigrations

import (
	"flag"
	"log"

	"gorm.io/gorm"
)

type ConfigDB struct {
	Driver   string
	Host     string
	Port     int
	Password string
	User     string
	Database string
}

var db *gorm.DB

func SetUp(config ConfigDB) {
	flag.Parse()

	args = setUpArgs()

	driver, err := getDB(config.Driver, config.Host, config.User, config.Password, config.Database, int16(config.Port))
	if err != nil {
		log.Fatal("Can't connec to database")
	}

	db = driver.getDB()
	getMigrations(db)
}

func Register(m func(*gorm.DB) IMigration) {

	exists := false
	migration := m(db)
	for _, v := range migrationsFiles {
		if v == migration {
			exists = true
		}
	}

	if exists {
		log.Fatalf("The migration %s has been register", migration.Name())
	}
	migrationsFiles = append(migrationsFiles, migration)
}

func Run() {
	command := getCommand(args, db)
	command.execute()
}
