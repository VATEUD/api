package subdivision

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Subdivisions(w http.ResponseWriter, r *http.Request) {
	var subdivisions []models.Subdivision

	if err := database.DB.Order("name asc").Find(&subdivisions).Error; err != nil {
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

func Subdivision(w http.ResponseWriter, r *http.Request) {
	var subdivision models.Subdivision

	attrs := mux.Vars(r)

	if err := database.DB.Where("code = ?", attrs["subdivision"]).Or("id = ?", attrs["subdivision"]).First(&subdivision).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Subdivision %s not found.\n", attrs["subdivision"])
			res := response.New(w, r, "Subdivision not found.", http.StatusNotFound)
			res.Process()
			return
		}

		log.Println("Error occurred while fetching subdivisions from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivisions.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(subdivision)

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
