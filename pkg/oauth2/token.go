package oauth2

import (
	"api/pkg/response"
	"log"
	"net/http"
)

func Token(w http.ResponseWriter, r *http.Request) {
	req, err := newAccessTokenRequest(r)

	if err != nil {
		log.Println("Error occurred while creating a new request. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
		res.Process()
		return
	}

	log.Println(req)
}
