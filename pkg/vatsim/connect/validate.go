package connect

import (
	"api/internal/pkg/database"
	"api/internal/pkg/logger"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	grantType = "authorization_code"
)

func Validate(w http.ResponseWriter, r *http.Request) {
	if err := r.URL.Query().Get("error"); len(err) > 0 {
		e := getError(err)

		res := response.New(w, r, e, http.StatusExpectationFailed)
		res.Process()
		return
	}

	code := r.URL.Query().Get("code")

	if len(code) < 0 {
		res := response.New(w, r, "Code was not provided.", http.StatusBadRequest)
		res.Process()
		return
	}

	tokenChannel := make(chan Token)
	go getToken(code, tokenChannel)
	token := <-tokenChannel

	if token.err != nil {
		logger.Log.Errorf("Internal server error occurred while fetching the user details 44. Error: %s.", token.err.Error())
		res := response.New(w, r, "Internal server error occurred while fetching the user details 44.", http.StatusInternalServerError)
		res.Process()
		return
	}

	userChannel := make(chan UserData)
	go getUserDetails(token, userChannel)
	user := <-userChannel

	if user.err != nil {
		logger.Log.Errorf("Internal server error occurred while fetching the user details 55. Error: %s.", user.err.Error())
		res := response.New(w, r, "Internal server error occurred while fetching the user details 55.", http.StatusInternalServerError)
		res.Process()
		return
	}

	saveChannel := make(chan error)
	go saveUser(user.Data, saveChannel)

	if err := <-saveChannel; err != nil {
		logger.Log.Errorf("Internal server error occurred while fetching the user details 66. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error occurred while fetching the user details 66.", http.StatusInternalServerError)
		res.Process()
		return
	}

	claims := jwt.MapClaims{
		"cid": user.Data.CID,
	}

	jwtToken, err := utils.GenerateNewJWT(claims)

	if err != nil {
		logger.Log.Errorf("Internal server error occurred while fetching the user details 79. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error occurred while fetching the user details 79.", http.StatusInternalServerError)
		res.Process()
		return
	}

	data := map[string]string{
		"token": jwtToken,
	}

	res, err := json.Marshal(data)

	if err != nil {
		logger.Log.Errorf("Internal server error occurred while fetching the user details 92. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error occurred while fetching the user details 92.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		logger.Log.Errorf("Internal server error occurred while fetching the user details 102. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error occurred while fetching the user details 102.", http.StatusInternalServerError)
		res.Process()
		return
	}
}

func getToken(code string, tokenChannel chan Token) {
	connectURL := utils.ConnectURL("oauth/token", "")
	data := url.Values{}
	data.Set("grant_type", grantType)
	data.Set("code", code)
	data.Set("redirect_uri", utils.Getenv("CONNECT_REDIRECT", ""))
	data.Set("client_id", utils.Getenv("CONNECT_CLIENT_ID", ""))
	data.Set("client_secret", utils.Getenv("CONNECT_CLIENT_SECRET", ""))

	client := &http.Client{}
	r, err := http.NewRequest("POST", connectURL.String(), strings.NewReader(data.Encode()))

	if err != nil {
		tokenChannel <- Token{err: err}
		return
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(r)

	if err != nil {
		tokenChannel <- Token{err: err}
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Failed to close the body. Error:", err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {
		tokenChannel <- Token{err: errors.New(fmt.Sprintf("expected code 200, got %d", res.StatusCode))}
		return
	}

	var token Token

	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		tokenChannel <- Token{err: err}
		return
	}

	tokenChannel <- token
}

func getUserDetails(token Token, userChannel chan UserData) {
	client := &http.Client{}
	connectURL := utils.ConnectURL("/api/user", "")
	req, err := http.NewRequest("GET", connectURL.String(), nil)

	if err != nil {
		userChannel <- UserData{err: err}
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.Access))

	res, err := client.Do(req)

	if err != nil {
		userChannel <- UserData{err: err}
		return
	}

	defer func(closer io.ReadCloser) {
		if err := closer.Close(); err != nil {
			log.Println("failed to close the body. Error:", err.Error())
		}
	}(res.Body)

	if res.StatusCode != 200 {

		if res.StatusCode == http.StatusUnauthorized {
			var errorResponse struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}

			if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
				userChannel <- UserData{err: errors.New(fmt.Sprintf("expected code 200, got %d. Error: %s.", res.StatusCode, err.Error()))}
				return
			}

			userChannel <- UserData{err: errors.New(errorResponse.Message)}
			return
		}

		userChannel <- UserData{err: errors.New(fmt.Sprintf("expected code 200, got %d", res.StatusCode))}
		return
	}

	var user UserData

	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		userChannel <- UserData{err: err}
		return
	}

	userChannel <- user
}

func saveUser(data Data, saveChannel chan error) {
	user := models.User{}
	if err := database.DB.Where("id = ?", data.CID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cid, err := strconv.Atoi(data.CID)
			if err != nil {
				saveChannel <- err
				return
			}

			user = models.User{
				ID:              uint(cid),
				NameFirst:       data.Personal.NameFirst,
				NameLast:        data.Personal.NameLast,
				Email:           data.Personal.Email,
				Rating:          data.Vatsim.Rating.ID,
				PilotRating:     data.Vatsim.PilotRating.ID,
				CountryID:       data.Personal.Country.ID,
				CountryName:     data.Personal.Country.Name,
				RegionID:        data.Vatsim.Region.ID,
				RegionName:      data.Vatsim.Region.Name,
				DivisionID:      data.Vatsim.Division.ID,
				DivisionName:    data.Vatsim.Division.Name,
				SubdivisionID:   data.Vatsim.Subdivision.ID,
				SubdivisionName: data.Vatsim.Subdivision.Name,
				CreatedAt:       time.Now().UTC(),
				UpdatedAt:       time.Now().UTC(),
			}

			if err := database.DB.Create(&user).Error; err != nil {
				saveChannel <- err
				return
			}
		} else {
			saveChannel <- err
			return
		}
	} else {
		user.NameFirst = data.Personal.NameFirst
		user.NameLast = data.Personal.NameLast
		user.Email = data.Personal.Email
		user.Rating = data.Vatsim.Rating.ID
		user.PilotRating = data.Vatsim.PilotRating.ID
		user.CountryID = data.Personal.Country.ID
		user.CountryName = data.Personal.Country.Name
		user.RegionID = data.Vatsim.Region.ID
		user.RegionName = data.Vatsim.Region.Name
		user.DivisionID = data.Vatsim.Division.ID
		user.DivisionName = data.Vatsim.Division.Name
		user.SubdivisionID = data.Vatsim.Subdivision.ID
		user.SubdivisionName = data.Vatsim.Subdivision.Name
		user.UpdatedAt = time.Now().UTC()

		if err := database.DB.Save(&user).Error; err != nil {
			saveChannel <- err
			return
		}
	}

	saveChannel <- nil
}

func getError(err string) string {
	authErrors := map[string]string{
		"invalid_request":           "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed.",
		"unauthorized_client":       "The client is not authorized to request an authorization code using this method.",
		"access_denied":             "The resource owner or authorization server denied the request.",
		"unsupported_response_type": "The authorization server does not support obtaining an authorization code using this method.",
		"invalid_scope":             "The requested scope is invalid, unknown, or malformed.",
		"server_error":              "The authorization server encountered an unexpected condition that prevented it from fulfilling the request.",
		"temporarily_unavailable":   "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server.",
	}

	e, ok := authErrors[err]

	if !ok {
		return "Error occurred."
	}

	return e
}
