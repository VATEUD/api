package division

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"net/http"
)

func Instructors(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")

	var instructors []models.DivisionInstructor
	if err := database.DB.Order("user_id asc").Find(&instructors).Error; err != nil {
		logger.Log.Errorln("Error occurred while fetching users from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching examiners.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(instructors)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching examiners.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
