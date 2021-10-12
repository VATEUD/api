package oauth2

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"fmt"
	"github.com/matoous/go-nanoid/v2"
	"log"
	"net/http"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	req, err := newRequest(r)

	if err != nil {
		log.Println("Error occurred while creating a new request. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
		res.Process()
		return
	}

	client, err := req.Validate()

	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, fmt.Sprintf("%s?%s", r.Form.Get("redirect_uri"), err.Error()), http.StatusFound)
		return
	}

	id, err := gonanoid.New(100)

	if err != nil {
		log.Println("Error occurred while generating the UID. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
		res.Process()
		return
	}

	login := models.OauthAuthCode{
		ID:        id,
		ClientID:  client.ID,
		Client:    *client,
		Scopes:    req.Scopes,
		UserAgent: r.UserAgent(),
	}

	if err := database.DB.Create(&login).Error; err != nil {
		log.Println("Error occurred while saving the auth code. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
		res.Process()
		return
	}
}
