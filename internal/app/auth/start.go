package auth

import (
	"auth/internal/pkg/database"
	"auth/pkg/web"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Start() {

	if _, err := os.Stat(".env"); err != nil {
		log.Fatalln("Environment file couldn't be found.")
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Failed to load the environment variables.")
		return
	}

	database.Connect()

	server := web.New()

	server.Start()
}
