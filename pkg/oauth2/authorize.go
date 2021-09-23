package oauth2

import (
	"auth/pkg/response"
	"log"
	"net/http"
	"text/template"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/template/authorize.html")

	if err != nil {
		log.Println("Error occurred while parsing the file. Error:", err.Error())
		res := response.New(w, r, "Internal server error occurred.", http.StatusInternalServerError)
		res.Process()
		return
	}

	tmpl.Execute(w, nil)
}
