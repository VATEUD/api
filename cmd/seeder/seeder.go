package main

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
	"api/utils"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
)

const (
	divisionCode = "EUD"
)

func main() {
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

	r := utils.NewRequest(fmt.Sprintf("%s/subdivisions", utils.Getenv("VATSIM_API_URL", "")), "GET", nil)

	success := make(chan bool)

	go r.Do(success)

	defer close(success)

	if true != <-success {
		logger.Log.Fatalln(r.Error.Error())
		return
	}

	defer r.Response.Body.Close()

	var subdivisions []models.Subdivision
	var subdivs []struct {
		Code           string `json:"code"`
		FullName       string `json:"fullname"`
		ParentDivision string `json:"parentdivision"`
	}

	body, err := ioutil.ReadAll(r.Response.Body)

	if err != nil {
		logger.Log.Fatalln(err.Error())
		return
	}

	if err := json.Unmarshal(body, &subdivs); err != nil {
		logger.Log.Fatalln(err.Error())
		return
	}

	for _, subdiv := range subdivs {
		s := models.Subdivision{
			Code: subdiv.Code,
			Name: subdiv.FullName,
		}

		if subdiv.ParentDivision == divisionCode {
			s.WithinVATEUD = true
		}

		subdivisions = append(subdivisions, s)
	}

	if err := database.DB.Create(subdivisions).Error; err != nil {
		logger.Log.Fatalln(err.Error())
		return
	}

	if err := database.DB.Exec("UPDATE subdivisions SET website_url = NULL").Error; err != nil {
		logger.Log.Fatalln(err.Error())
	}
}
