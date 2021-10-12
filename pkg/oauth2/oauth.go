package oauth2

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"api/pkg/vatsim/connect"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func User(w http.ResponseWriter, r *http.Request) {
	cid := r.Header.Get("cid")
	user := models.User{}
	if err := database.DB.Where("id = ?", cid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found. CID #%s.\n", cid)
			res := response.New(w, r, fmt.Sprintf("User not found. CID #%d.\n", cid), http.StatusNotFound)
			res.Process()
			return
		}

		log.Printf("Error occurred while executing the query on /api/user. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error while fetching user information.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := connectJson(user, []string{"full_name"})

	if err != nil {
		log.Printf("Error occurred while marshalling the response on /api/user. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error while fetching user information.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if _, err := w.Write(bytes); err != nil {
		log.Println(err.Error())
	}
}

func connectJson(user models.User, scopes []string) ([]byte, error) {
	cUser := connect.UserData{}

	cUser.Data.CID = fmt.Sprintf("%d", user.ID)

	if isInScopes(scopes, "full_name") {
		cUser.Data.Personal.NameFirst = user.NameFirst
		cUser.Data.Personal.NameLast = user.NameLast
		cUser.Data.Personal.NameFull = fmt.Sprintf("%s %s", user.NameFirst, user.NameLast)
	}

	if isInScopes(scopes, "country") {
		cUser.Data.Personal.Country = connect.Country{
			ID:   user.CountryID,
			Name: user.CountryName,
		}
	}

	if isInScopes(scopes, "email") {
		cUser.Data.Personal.Email = user.Email
	}

	if isInScopes(scopes, "vatsim_details") {
		cUser.Data.Vatsim = connect.Vatsim{
			Rating: connect.Rating{
				ID: user.Rating,
			},
			PilotRating: connect.PilotRating{
				ID: user.PilotRating,
			},
			Region: connect.Region{
				ID:   user.RegionID,
				Name: user.RegionName,
			},
			Division: connect.Division{
				ID:   user.DivisionID,
				Name: user.DivisionName,
			},
			Subdivision: connect.Subdivision{
				ID:   user.SubdivisionID,
				Name: user.SubdivisionName,
			},
		}
	}

	return json.Marshal(cUser)
}

func isInScopes(scopes []string, scope string) bool {
	for _, s := range scopes {
		if scope == s {
			return true
		}
	}

	return false
}
