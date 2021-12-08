package gormmigrations

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type idb interface {
	getDB() *gorm.DB
}

type database struct {
	host     string
	password string
	user     string
	name     string
	port     int16
	*gorm.DB
}

type postgresDB struct {
	database
}

func newPostgresDB(host string, user string, password string, db string, port int16) idb {
	return &postgresDB{
		database: database{
			host:     host,
			user:     user,
			password: password,
			port:     port,
			name:     db,
		},
	}
}

func (db *postgresDB) getDB() *gorm.DB {
	if db.DB == nil {
		url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable connect_timeout=5",
			db.host,
			db.user,
			db.password,
			db.name,
			db.port,
		)

		DB, err := gorm.Open(postgres.Open(url), &gorm.Config{})

		if err != nil {
			fmt.Println("Can't connect to database")
			panic(err)
		}

		db.DB = DB
	}

	return db.DB
}

type mysqlDB struct {
	database
}

func newMysqlDB(host string, user string, password string, db string, port int16) idb {
	return &mysqlDB{
		database: database{
			host:     host,
			user:     user,
			password: password,
			port:     port,
			name:     db,
		},
	}
}

func (db *mysqlDB) getDB() *gorm.DB {
	if db.DB == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.user,
			db.password,
			db.host,
			db.port,
			db.name,
		)

		DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			fmt.Println("Can't connect to database")
			panic(err)
		}

		db.DB = DB
	}

	return db.DB
}

func getDB(driver string, host string, user string, password string, db string, port int16) (idb, error) {
	switch driver {
	case "postgres":
		return newPostgresDB(host, user, password, db, port), nil
	case "mysql":
		return newMysqlDB(host, user, password, db, port), nil
	default:
		return nil, errors.New("Driver not supported")
	}
}
