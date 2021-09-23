package database

import (
	"api/utils"
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxAttempts = 3
)

var (
	DB      *Database
	attempt = 1
)

func Connect() {
	DB = &Database{
		API:     connect(utils.Getenv("API_DB_NAME", "")),
		Central: connect(utils.Getenv("CENTRAL_DB_NAME", "")),
	}
}

func connect(database string) *gorm.DB {
	config := retrieveDatabaseCredentials(database)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", config.user, config.password, config.hostname, config.port, config.database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		if attempt <= maxAttempts {
			log.Errorf("Error connecting to the database. Error: %s.", err.Error())
			attempt += 1
			Connect()
			return nil
		}

		panic("Failed to connect to the database. Aborting...")
	}

	return db
}
