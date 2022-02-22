package vatsim

import (
	"api/internal/pkg/logger"
	"api/pkg/response"
	"api/pkg/vatsim/webhook"
	"api/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	authHeaderPrefix = "Token"
	userAgent        = "VATSIM-API"
	contenType       = "application/json"
)

func Webhook(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")

	if contentType != contenType {
		logger.Log.Println("Unauthenticated")
		res := response.New(w, r, "Unauthorized", http.StatusUnauthorized)
		res.Process()
		return
	}

	userAgentHeader := r.Header.Get("User-Agent")

	if userAgentHeader != userAgent {
		logger.Log.Println("Unauthenticated")
		res := response.New(w, r, "Unauthorized", http.StatusUnauthorized)
		res.Process()
		return
	}

	authHeader := strings.TrimPrefix(r.Header.Get("Authorization"), authHeaderPrefix)

	if len(authHeader) != 35 && authHeader != utils.Getenv("VATSIM_WEBHOOK_TOKEN", "") {
		logger.Log.Println("Unauthenticated")
		res := response.New(w, r, "Unauthorized", http.StatusUnauthorized)
		res.Process()
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Log.Errorln("Error occurred while reading the request body. Error:", err.Error())
		res := response.New(w, r, "Internal server error while reading the request body.", http.StatusInternalServerError)
		res.Process()
		return
	}

	var event webhook.Member

	if err := json.Unmarshal(body, &event); err != nil {
		logger.Log.Errorln("Error occurred while unmarshalling the request body. Error:", err.Error())
		res := response.New(w, r, "Internal server error while reading the request body.", http.StatusInternalServerError)
		res.Process()
		return
	}

	resultChannel := make(chan webhook.Result)

	if event.Type == "member_changed_webhook" {
		go event.Update(resultChannel)
	}

	res := <-resultChannel

	close(resultChannel)

	resp := response.New(w, r, res.Message, res.StatusCode)
	resp.Process()
}
