package news

import (
	"api/internal/pkg/database"
	"api/pkg/models/api"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"log"
	"net/http"
)

func News(w http.ResponseWriter, r *http.Request) {
	var news []api.News

	if err := database.DB.API.Find(&news).Error; err != nil {
		log.Println("Error occurred while fetching news from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching news.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(news)

	if err != nil {
		log.Println("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching news.", http.StatusInternalServerError)
		res.Process()
		return
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
