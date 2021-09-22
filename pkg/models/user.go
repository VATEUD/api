package models

import (
	"auth/pkg/vatsim/connect"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey;type:bigint(20);unsigned;column:id"`
	NameFirst       string    `json:"name_first" gorm:"type:varchar(255);column:name_first"`
	NameLast        string    `json:"name_last" gorm:"type:varchar(255);column:name_last"`
	Email           string    `json:"-" gorm:"type:varchar(255);column:email"`
	Rating          int       `json:"rating" gorm:"type:tinyint(4);column:rating"`
	PilotRating     int       `json:"pilot_rating" gorm:"type:tinyint(4);column:pilot_rating"`
	CountryID       string    `json:"country_id" gorm:"type:varchar(255);column:country_id"`
	CountryName     string    `json:"country_name" gorm:"type:varchar(255);column:country_name"`
	RegionID        string    `json:"region_id" gorm:"type:varchar(255);column:region_id"`
	RegionName      string    `json:"region_name" gorm:"type:varchar(255);column:region_name"`
	DivisionID      string    `json:"division_id" gorm:"type:varchar(255);column:division_id"`
	DivisionName    string    `json:"division_name" gorm:"type:varchar(255);column:division_name"`
	SubdivisionID   string    `json:"subdivision_id" gorm:"type:varchar(255);column:subdivision_id"`
	SubdivisionName string    `json:"subdivision_name" gorm:"type:varchar(255);column:subdivision_name"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:timestamp;column:updated_at"`
}

func (user User) Json() ([]byte, error) {
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
