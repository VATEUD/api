package news

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

func NewsIndex(w http.ResponseWriter, r *http.Request) {
	var news []models.News

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
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing response. Error:", err.Error())
	}
}

func NewsShow(w http.ResponseWriter, r *http.Request) {
	attrs := mux.Vars(r)

	news := models.News{}

	if err := database.DB.API.Where("id = ?", attrs["id"]).First(&news).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("News #%s not found.\n", attrs["id"])
			res := response.New(w, r, "News article not found.", http.StatusNotFound)
			res.Process()
			return
		}

		log.Println("Error occurred while fetching news article from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching the news article.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := json.Marshal(news)

	if err != nil {
		log.Println("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching the news article.", http.StatusInternalServerError)
		res.Process()
		return
	}

	utils.Allow(w, "*")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing the response. Error:", err.Error())
	}
}
