package oauth2

import (
	"api/pkg/response"
	"log"
	"net/http"
)

func Token(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	req, err := newAccessTokenRequest(r)

	if err != nil {
		log.Println("Error occurred while creating a new request. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
		res.Process()
		return
	}

	_, aErr := req.Validate()

	if aErr != nil {
		log.Println("Validation failed. Error:", aErr.internalError.Error())
		data, err := aErr.Json()

		if err != nil {
			log.Println("Error occurred while marshalling the json response. Error:", err.Error())
			res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
			res.Process()
			return
		}

		w.WriteHeader(aErr.Code)
		if _, err := w.Write(data); err != nil {
			log.Println("failed to write")
			return
		}
	}
}
