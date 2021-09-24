package subdivision

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"log"
	"net/http"
)

func Subdivisions(w http.ResponseWriter, r *http.Request) {
	var subdivisions []models.Subdivision

	if err := database.DB.Find(&subdivisions).Error; err != nil {
		log.Println("Error occurred while fetching subdivisions from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivisions.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(subdivisions)

	if err != nil {
		log.Println("Error occurred while marshalling data. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivisions.", http.StatusInternalServerError)
		res.Process()
		return
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error occurred while writing the response. Error:", err.Error())
	}
}