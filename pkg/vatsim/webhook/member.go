package webhook

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
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
		logger.Log.Println("User not found")
		result <- Result{
			StatusCode: http.StatusOK,
			Message:    "Member updated",
		}
		return
	}

	user.NameFirst = member.NameFirst
	user.NameLast = member.NameLast
	user.Rating = member.Rating
	user.PilotRating = member.PilotRating
	user.RegionID = member.RegionId
	user.DivisionID = member.DivisionId
	user.SubdivisionID = member.SubdivisionId

	if err := database.DB.Updates(&user).Error; err != nil {
		logger.Log.Errorln("Failed to update member's data. Error:", err.Error())
		result <- Result{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to update member's data",
		}
		return
	}

	result <- Result{
		StatusCode: http.StatusOK,
		Message:    "Member updated",
	}
}
