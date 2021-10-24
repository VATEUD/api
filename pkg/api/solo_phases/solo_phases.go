package solo_phases

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func RetrieveAll(w http.ResponseWriter, r *http.Request) {
	var soloPhases []*models.SoloPhase

	if err := database.DB.Where("expired = ?", false).Preload("Subdivision").Find(&soloPhases).Error; err != nil {
		log.Println("Error occurred while fetching solo phases from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	b, err := json.Marshal(soloPhases)

	if err != nil {
		log.Println("Error occurred while marshalling the solo phases. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Println("Failed to write the response. Error:", err.Error())
	}
}

func RetrieveBySubdivision(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var subdivision models.Subdivision
	var solos []*models.SoloPhase

	if err := database.DB.Where("code = ?", params["subdivision"]).Or("id = ?", params["subdivision"]).First(&subdivision).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Subdivision not found. Error:", err.Error())
			res := response.New(w, r, "Subdivision you are looking for, couldn't be found.", http.StatusNotFound)
			res.Process()
			return
		}

		log.Println("Error occurred while fetching solo phases from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := database.DB.Where("subdivision_id = ?", subdivision.ID).Find(&solos).Error; err != nil {
		log.Println("Error occurred while fetching solo phases from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	b, err := json.Marshal(solos)

	if err != nil {
		log.Println("Error occurred while marshalling the solo phases. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Println("Failed to write the response. Error:", err.Error())
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	req, err := newSoloPhaseRequest(r)

	if err != nil {
		if errors.Is(err, userIDNotProvided) || errors.Is(err, validUntilNotProvided) {
			log.Println("Error occurred while creating a new solo phase request. Error:", err.Error())
			res := response.New(w, r, fmt.Sprintf("Make sure all necessary fields are provided. %s.", err.Error()), http.StatusInternalServerError)
			res.Process()
			return
		}

		log.Println("Error occurred while creating a new solo phase request. Error:", err.Error())
		res := response.New(w, r, "Internal server error while creating the solo phase.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := validate(req); err != nil {
		log.Println("Validation failed. Error:", err.Error())
		res := response.New(w, r, fmt.Sprintf("Validation failed. Error: %s.", err.Error()), http.StatusInternalServerError)
		res.Process()
		return
	}

	check := database.DB.Where("user_id = ? AND expired = ?", req.User, false).Find(&models.SoloPhase{})

	if check.Error == nil {
		if check.RowsAffected > 0 {
			log.Println("user already has an active solo phase")
			res := response.New(w, r, "Given member already has an active solo phase", http.StatusForbidden)
			res.Process()
			return
		}
	}

	token := r.Context().Value("token").(models.SubdivisionToken)

	saveChannel := make(chan soloPhaseSaveResult)
	go req.save(saveChannel, token.Subdivision)

	solo := <- saveChannel

	b, err := json.Marshal(solo.soloPhase)

	if err != nil {
		log.Println("Error occurred while marshalling the solo phase. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred while creating the solo phase", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Println("error writing the response")
	}
}
