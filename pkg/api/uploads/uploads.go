package uploads

import (
	"api/internal/pkg/database"
	"api/pkg/minio"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")

	var uploads []*models.Upload

	if err := database.DB.Where("public = ?", true).Find(&uploads).Error; err != nil {
		log.Println("Error occurred while fetching uploads from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching uploads.", http.StatusInternalServerError)
		res.Process()
		return
	}

	for _, upload := range uploads {
		upload.SetDownloadURL()
	}

	bytes, err := json.Marshal(uploads)

	if err != nil {
		log.Println("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching uploads.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing response. Error:", err.Error())
	}
}

func Filter(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")

	var uploads []*models.Upload

	attrs := mux.Vars(r)

	if err := database.DB.Where("public = ? AND type = ?", true, attrs["type"]).Find(&uploads).Error; err != nil {
		log.Println("Error occurred while fetching uploads from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching uploads.", http.StatusInternalServerError)
		res.Process()
		return
	}

	for _, upload := range uploads {
		upload.SetDownloadURL()
	}

	bytes, err := json.Marshal(uploads)

	if err != nil {
		log.Println("Error occurred while marshalling the response. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching uploads.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); err != nil {
		log.Println("Error writing response. Error:", err.Error())
	}
}

func Download(w http.ResponseWriter, r *http.Request) {
	utils.Allow(w, "*")

	attrs := mux.Vars(r)

	upload := models.Upload{}

	if err := database.DB.Where("public = ? AND id = ?", true, attrs["id"]).Find(&upload).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Upload not found. Error:", err.Error())
			res := response.New(w, r, "File you are looking for, could not be found.", http.StatusNotFound)
			res.Process()
		}

		log.Println("Error occurred while fetching upload from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching the file.", http.StatusInternalServerError)
		res.Process()
		return
	}

	client, err := minio.New()

	if err != nil {
		log.Println("Error occurred while starting the minio session. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching the file.", http.StatusInternalServerError)
		res.Process()
		return
	}

	file, err := client.Download(upload.Path)

	if err != nil {
		log.Println("Error occurred while fetching the file from storage. Error:", err.Error())
		res := response.New(w, r, "File not found.", http.StatusNotFound)
		res.Process()
		return
	}

	w.Header().Set("Content-Type", *file.Object.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	io.Copy(w, file.Object.Body)
}
