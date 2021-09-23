package database

import "gorm.io/gorm"

type credentials struct {
	user, password, hostname, port, database string
}

type Database struct {
	API, Central *gorm.DB
}
