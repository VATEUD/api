package database

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	maxAttempts = 3
)

var (
	DB      *gorm.DB
	attempt = 1
)

func Connect() {
	config := retrieveDatabaseCredentials()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", config.user, config.password, config.hostname, config.port, config.database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		if attempt <= maxAttempts {
			log.Errorf("Error connecting to the database. Error: %s.", err.Error())
			attempt += 1
			Connect()
			return
		}

		panic("Failed to connect to the database. Aborting...")
	}

	DB = db
}
