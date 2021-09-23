package oauth2

import (
	"auth/internal/pkg/database"
	"auth/pkg/models/central"
	"auth/pkg/response"
	"auth/pkg/vatsim/connect"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func User(w http.ResponseWriter, r *http.Request) {
	cid := r.Header.Get("cid")
	user := central.User{}
	if err := database.DB.Central.Where("id = ?", cid).First(&user).Error; err != nil {
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

	bytes, err := connectJson(user)

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

func connectJson(user central.User) ([]byte, error) {
	res := connect.UserData{Data: connect.Data{
		CID: fmt.Sprintf("%d", user.ID),
		Personal: connect.Personal{
			NameFirst: user.NameFirst,
			NameLast:  user.NameLast,
			NameFull:  fmt.Sprintf("%s %s", user.NameFirst, user.NameLast),
			Email:     user.Email,
			Country: connect.Country{
				ID:   user.CountryID,
				Name: user.CountryName,
			},
		},
		Vatsim: connect.Vatsim{
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
		},
	}}

	return json.Marshal(res)
}
