package solo_phases

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
	"api/pkg/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

func RetrieveAll(w http.ResponseWriter, r *http.Request) {
	var soloPhases []*models.SoloPhase

	if err := database.DB.Where("expired = ?", false).Preload("Subdivision").Find(&soloPhases).Error; err != nil {
		logger.Log.Errorln("Error occurred while fetching solo phases from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	b, err := json.Marshal(soloPhases)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling the solo phases. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		logger.Log.Errorln("Failed to write the response. Error:", err.Error())
	}
}

func RetrieveBySubdivision(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var subdivision models.Subdivision
	var solos []*models.SoloPhase

	if err := database.DB.Where("code = ?", params["subdivision"]).Or("id = ?", params["subdivision"]).First(&subdivision).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Println("Subdivision not found. Error:", err.Error())
			res := response.New(w, r, "Subdivision you are looking for, couldn't be found.", http.StatusNotFound)
			res.Process()
			return
		}

		logger.Log.Errorln("Error occurred while fetching solo phases from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := database.DB.Where("subdivision_id = ?", subdivision.ID).Find(&solos).Error; err != nil {
		logger.Log.Errorln("Error occurred while fetching solo phases from the DB. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	b, err := json.Marshal(solos)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling the solo phases. Error:", err.Error())
		res := response.New(w, r, "Internal server error while fetching solo phases.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		logger.Log.Errorln("Failed to write the response. Error:", err.Error())
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	req, err := newSoloPhaseRequest(r)

	if err != nil {
		if errors.Is(err, userIDNotProvided) || errors.Is(err, validUntilNotProvided) {
			logger.Log.Println("Error occurred while creating a new solo phase request. Error:", err.Error())
			res := response.New(w, r, fmt.Sprintf("Make sure all necessary fields are provided. %s.", err.Error()), http.StatusInternalServerError)
			res.Process()
			return
		}

		logger.Log.Errorln("Error occurred while creating a new solo phase request. Error:", err.Error())
		res := response.New(w, r, "Internal server error while creating the solo phase.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := validate(req); err != nil {
		logger.Log.Errorln("Validation failed. Error:", err.Error())
		res := response.New(w, r, fmt.Sprintf("Validation failed. Error: %s.", err.Error()), http.StatusInternalServerError)
		res.Process()
		return
	}

	check := database.DB.Where("user_id = ? AND expired = ?", req.User, false).Find(&models.SoloPhase{})

	if check.Error == nil {
		if check.RowsAffected > 0 {
			logger.Log.Println("user already has an active solo phase")
			res := response.New(w, r, "Given member already has an active solo phase", http.StatusForbidden)
			res.Process()
			return
		}
	}

	token := r.Context().Value("token").(models.SubdivisionToken)

	saveChannel := make(chan soloPhaseSaveResult)
	go req.save(saveChannel, token.Subdivision)

	solo := <-saveChannel

	if solo.err != nil {
		logger.Log.Errorln("Error occurred while creating the solo phase. Error:", solo.err.Error())
		res := response.New(w, r, "Internal server error occurred while creating the solo phase", http.StatusInternalServerError)
		res.Process()
		return
	}

	b, err := json.Marshal(solo.soloPhase)

	if err != nil {
		logger.Log.Errorln("Error occurred while marshalling the solo phase. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred while creating the solo phase", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		logger.Log.Errorln("error writing the response")
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	token := r.Context().Value("token").(models.SubdivisionToken)

	var solo models.SoloPhase

	if err := database.DB.Where("id = ? AND subdivision_id = ? AND expired = ?", params["id"], token.SubdivisionID, false).First(&solo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Println("Solo phase not found. ID:", params["id"])
			res := response.New(w, r, "Solo phase not found.", http.StatusNotFound)
			res.Process()
			return
		}

		logger.Log.Errorln("Error occurred while fetching the solo phase. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred while deleting the solo phase.", http.StatusInternalServerError)
		res.Process()
		return
	}

	go func() {
		solo.Expired = true
		database.DB.Updates(&solo)
	}()

	res := response.New(w, r, "Solo phase deleted", http.StatusOK)
	res.Process()
}

func Extend(w http.ResponseWriter, r *http.Request) {

	req, err := newExtendSoloPhaseRequest(r)

	if err != nil {
		if errors.Is(err, amountNotProvided) {
			logger.Log.Println(err.Error())
			res := response.New(w, r, "Amount was not provided", http.StatusBadRequest)
			res.Process()
			return
		}

		logger.Log.Errorln("Failed to create extend solo request. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred while extending the solo phase.", http.StatusInternalServerError)
		res.Process()
		return
	}

	if err := validate(req); err != nil {
		logger.Log.Errorln("Validation failed. Error:", err.Error())
		res := response.New(w, r, fmt.Sprintf("Validation failed. Error: %s.", err.Error()), http.StatusBadRequest)
		res.Process()
		return
	}

	params := mux.Vars(r)

	token := r.Context().Value("token").(models.SubdivisionToken)

	var solo models.SoloPhase

	if err := database.DB.Where("id = ? AND subdivision_id = ? AND expired = ?", params["id"], token.SubdivisionID, false).First(&solo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Println("Solo phase not found. ID:", params["id"])
			res := response.New(w, r, "Solo phase not found.", http.StatusNotFound)
			res.Process()
			return
		}

		logger.Log.Errorln("Error occurred while fetching the solo phase. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred while deleting the solo phase.", http.StatusInternalServerError)
		res.Process()
		return
	}

	saveChannel := make(chan soloPhaseSaveResult)
	go req.save(saveChannel, solo)

	save := <-saveChannel

	if save.err != nil {
		logger.Log.Errorln("Error occurred while saving the solo phase. Error:", save.err.Error())
		res := response.New(w, r, "Internal server error occurred while extending the solo phase.", http.StatusInternalServerError)
		res.Process()
		return
	}

	res := response.New(w, r, "Solo phase extended.", http.StatusOK)
	res.Process()
}
