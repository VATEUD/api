package api

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/cache"
	"api/pkg/web"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Start() {

	log.Println("Checking if .env file exists")
	if _, err := os.Stat(".env"); err != nil {
		log.Fatalln("Environment file couldn't be found.")
		return
	}

	log.Println("Loading environment variables")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Failed to load the environment variables.")
		return
	}

	if err := logger.New(); err != nil {
		panic(err)
		return
	}

	logger.Log.Println("Connecting to the database")
	database.Connect()
	cache.New()

	server := web.New()

	logger.Log.Println("Starting the web server")
	if err := server.Start(); err != nil {
		logger.Log.Fatalln("Error starting the web server. Error:", err.Error())
	}
}
