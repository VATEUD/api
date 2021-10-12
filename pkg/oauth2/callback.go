package oauth2

import (
	"api/internal/pkg/database"
	"api/pkg/models"
	"api/pkg/response"
	"api/pkg/vatsim/connect"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	_, err := connect.Validate(w, r)

	if err != nil {
		log.Printf("Failed to fetch user details. Error: %s.\n", err.Error())
		res := response.New(w, r, "Failed to fetch user details.", http.StatusInternalServerError)
		res.Process()
		return
	}

	cookie, err := r.Cookie(cookieName)

	if err != nil {
		log.Printf("Failed to fetch the cookie. Error: %s.\n", err.Error())
		res := response.New(w, r, "Failed to retrieve the token.", http.StatusInternalServerError)
		res.Process()
		return
	}

	defer func() {
		cookie = &http.Cookie{
			Name:    cookieName,
			Value:   "",
			Path:    "/",
			Domain:  r.URL.Host,
			Expires: time.Now(),
		}

		http.SetCookie(w, cookie)
	}()

	authCode := models.OauthAuthCode{}
	if err := database.DB.Debug().Preload("Client").Where("id = ? AND created_at > ?", cookie.Value, time.Now().UTC().Add(-time.Minute*5)).First(&authCode).Error; err != nil {
		log.Printf("Failed to fetch the token. Error: %s.\n", err.Error())
		res := response.New(w, r, "Invalid token.", http.StatusInternalServerError)
		res.Process()
		return
	}

	params := url.Values{}
	params.Set("code", authCode.ID)

	if len(authCode.State) > 0 {
		params.Set("state", authCode.State)
	}

	http.Redirect(w, r, fmt.Sprintf("%s?%s", authCode.Client.Redirect, params.Encode()), http.StatusFound)
}
