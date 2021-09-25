package division

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"log"
	"net/http"
)

func Staff(w http.ResponseWriter, r *http.Request) {
	var departments []models.StaffDepartment

	if err := database.DB.Preload("Members").Find(&departments).Error; err != nil {
		log.Printf("Error occurred while executing the query. Error: %s.\n", err.Error())
		res := response.New(w, r, "Internal server while fetching staff members.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(departments)

	if err != nil {
		log.Printf("Error occurred marshalling the response. Error: %s.\n", err.Error())
		res := response.New(w, r, "Internal server while fetching staff members.", http.StatusInternalServerError)
		res.Process()
		return
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing response.")
	}
}
