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

func Staff(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")
	var departments []models.StaffDepartment

	if err := database.DB.Preload("Members").Find(&departments).Error; err != nil {
		logger.Log.Errorln("Error occurred while executing the query. Error: %s.\n", err.Error())
		res := response.New(w, r, "Internal server while fetching staff members.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(departments)

	if err != nil {
		logger.Log.Errorln("Error occurred marshalling the response. Error: %s.\n", err.Error())
		res := response.New(w, r, "Internal server while fetching staff members.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		logger.Log.Errorln("Error writing response.")
	}
}
