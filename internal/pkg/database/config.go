package database

import (
	"api/utils"
	"fmt"
	"strings"
)

const (
	defaultDatabaseUsername = "DB_USERNAME"
	defaultDatabasePassword = "DB_PASSWORD"
	defaultDatabaseHost     = "DB_HOST"
	defaultDatabasePort     = "DB_PORT"
)

// retrieveDatabaseCredentials retrieves credentials that are needed for authentication in order to execute mysqldump command
func retrieveDatabaseCredentials(db string) credentials {
	if !isUsingSingleCredentials() {
		name := strings.ToUpper(db)

		return credentials{
			user:     utils.Getenv(fmt.Sprintf("%s_%s", name, defaultDatabaseUsername), ""),
			password: utils.Getenv(fmt.Sprintf("%s_%s", name, defaultDatabasePassword), ""),
			hostname: utils.Getenv(fmt.Sprintf("%s_%s", name, defaultDatabaseHost), ""),
			port:     utils.Getenv(fmt.Sprintf("%s_%s", name, defaultDatabasePort), ""),
			database: db,
		}
	}

	return credentials{
		user:     utils.Getenv(defaultDatabaseUsername, ""),
		password: utils.Getenv(defaultDatabasePassword, ""),
		hostname: utils.Getenv(defaultDatabaseHost, ""),
		port:     utils.Getenv(defaultDatabasePort, ""),
		database: db,
	}
}

func isUsingSingleCredentials() bool {
	return utils.Getenv("SINGLE_DB_CREDENTIALS", "false") == "true"
}
