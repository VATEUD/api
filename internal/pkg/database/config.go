package database

import (
	"auth/utils"
)

const (
	defaultDatabaseUsername = "DB_USERNAME"
	defaultDatabasePassword = "DB_PASSWORD"
	defaultDatabaseHost     = "DB_HOST"
	defaultDatabasePort     = "DB_PORT"
	defaultDatabaseName     = "DB_NAME"
)

// retrieveDatabaseCredentials retrieves credentials that are needed for authentication in order to execute mysqldump command
func retrieveDatabaseCredentials() config {
	return config{
		user:     utils.Getenv(defaultDatabaseUsername, ""),
		password: utils.Getenv(defaultDatabasePassword, ""),
		hostname: utils.Getenv(defaultDatabaseHost, ""),
		port:     utils.Getenv(defaultDatabasePort, ""),
		database: utils.Getenv(defaultDatabaseName, ""),
	}
}
