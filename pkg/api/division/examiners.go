package division

import (
	"auth/internal/pkg/database"
	"auth/pkg/models/api"
	"auth/pkg/response"
	"auth/utils"
	"encoding/json"
	"log"
	"net/http"
)

func Examiners(w http.ResponseWriter, r *http.Request) {
	var examiners []*api.DivisionExaminer
	if err := database.DB.API.Find(&examiners).Error; err != nil {
		log.Println("Error occurred while fetching users from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching examiners.", http.StatusInternalServerError)
		res.Process()
		return
	}
	utils.Allow(w, "*")

	for _, examiner := range examiners {
		examiner.SetUpTo()
	}

	bytes, err := json.Marshal(examiners)

	if err != nil {
		log.Println("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching examiners.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}