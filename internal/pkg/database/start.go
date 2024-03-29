package database

import (
	"api/internal/pkg/logger"
	"api/utils"
	"fmt"
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
	DB = connect(utils.Getenv("CENTRAL_DB_NAME", ""))
}

func connect(database string) *gorm.DB {
	config := retrieveDatabaseCredentials(database)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", config.user, config.password, config.hostname, config.port, config.database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		if attempt <= maxAttempts {
			logger.Log.Printf("Error connecting to the database. Error: %s.", err.Error())
			attempt += 1
			Connect()
			return nil
		}

		logger.Log.Panicf("Error connecting to the database. Error: %s.", err.Error())
		panic("Failed to connect to the database. Aborting...")
	}

	return db
}
