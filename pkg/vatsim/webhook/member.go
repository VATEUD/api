package webhook

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type Member struct {
	Event
	Id               uint   `json:"id"`
	NameFirst        string `json:"name_first"`
	NameLast         string `json:"name_last"`
	Rating           int    `json:"rating"`
	PilotRating      int    `json:"pilotrating"`
	SuspDate         string `json:"susp_date"`
	RegDate          string `json:"reg_date"`
	RegionId         string `json:"region_id"`
	DivisionId       string `json:"division_id"`
	SubdivisionId    string `json:"subdivision_id"`
	LastRatingChange string `json:"last_rating_change"`
}

func (member Member) Update(result chan Result) {
	var user models.User

	if err := database.DB.Where("id = ?", member.Id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Println("User not found")
			result <- Result{
				StatusCode: http.StatusOK,
				Message:    "Member updated",
			}
			return
		}

		logger.Log.Errorln("Error occurred while executing the query. Error:", err.Error())
		result <- Result{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}
		return
	}

	var subdivision models.Subdivision

	subError := database.DB.Where("code = ?", member.SubdivisionId).First(&subdivision).Error

	if subError != nil && !errors.Is(subError, gorm.ErrRecordNotFound) {
		logger.Log.Errorln("Error occurred while processing the query. Error:", subError.Error())
		result <- Result{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}
		return
	}

	user.NameFirst = member.NameFirst
	user.NameLast = member.NameLast
	user.Rating = member.Rating
	user.PilotRating = member.PilotRating
	user.RegionID = member.RegionId
	user.RegionName = Regions[member.RegionId]
	user.DivisionID = member.DivisionId
	user.DivisionName = Divisions[member.DivisionId]
	user.SubdivisionID = subdivision.Code
	user.SubdivisionName = subdivision.Name

	if err := database.DB.Updates(&user).Error; err != nil {
		logger.Log.Errorln("Failed to update member's data. Error:", err.Error())
		result <- Result{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update member's data",
		}
		return
	}

	// TODO - investigate why subdivision doesn't get saved if it's set to null
	if errors.Is(subError, gorm.ErrRecordNotFound) {
		if err := database.DB.Exec("UPDATE users SET subdivision_id = NULL, subdivision_name = NULL WHERE id = ?", user.ID).Error; err != nil {
			logger.Log.Errorln("Failed to update member's data. Error:", err.Error())
			result <- Result{
				StatusCode: http.StatusInternalServerError,
				Message:    "Failed to update member's data",
			}
			return
		}

	}

	logger.Log.Println("Received member's data -", member)

	result <- Result{
		StatusCode: http.StatusOK,
		Message:    "Member updated",
	}
}
