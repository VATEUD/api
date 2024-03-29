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

func Examiners(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")

	var examiners []*models.DivisionExaminer
	if err := database.DB.Find(&examiners).Error; err != nil {
		logger.Log.Errorln("Error occurred while fetching users from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching examiners.", http.StatusInternalServerError)
		res.Process()
		return
	}

	for _, examiner := range examiners {
		examiner.SetUpTo()
	}

	bytes, err := json.Marshal(examiners)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching examiners.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
