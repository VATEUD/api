package oauth2

import (
	"api/pkg/response"
	"fmt"
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

	if err := req.Validate(); err != nil {
		http.Redirect(w, r, fmt.Sprintf("%s?%s", r.Form.Get("redirect_uri"), err.Error()), http.StatusFound)
		return
	}
}
