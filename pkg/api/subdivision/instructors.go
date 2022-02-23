package subdivision

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func Instructors(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")
	var instructors []models.Subdivision

	if err := database.DB.Preload("Instructors").Where("within_eud = ?", true).Order("name asc").Find(&instructors).Error; err != nil {
		logger.Log.Errorln("Error occurred while fetching subdivision instructors from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivision instructors.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(instructors)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling data. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivision instructors.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		logger.Log.Errorln("Error occurred while writing the response. Error:", err.Error())
	}
}

func InstructorsFilter(w http.ResponseWriter, r *http.Request) {
	var instructors []models.SubdivisionInstructor
	var subdivision models.Subdivision

	attrs := mux.Vars(r)

	if err := database.DB.Where("code = ? AND within_eud = ?", attrs["subdivision"], true).First(&subdivision).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Printf("Subdivision %s not found.\n", attrs["subdivision"])
			res := response.New(w, r, "Subdivision not found.", http.StatusNotFound)
			res.Process()
			return
		}

		logger.Log.Errorln("Error occurred while fetching subdivisions from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivision instructors.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := database.DB.Where("subdivision_id = ?", subdivision.ID).Find(&instructors).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Printf("Subdivision %s not found.\n", attrs["subdivision"])
			res := response.New(w, r, "Instructors not found.", http.StatusNotFound)
			res.Process()
			return
		}

		logger.Log.Errorln("Error occurred while fetching subdivisions from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivision instructors.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(instructors)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling data. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching subdivision instructors.", http.StatusInternalServerError)
		res.Process()
		return
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		logger.Log.Errorln("Error occurred while writing the response. Error:", err.Error())
	}
}
