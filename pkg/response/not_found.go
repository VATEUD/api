package response

import (
	"log"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Page %s could not be found.\n", r.URL.String())
	res := New(w, r, "The page you are looking for could not be found.", http.StatusNotFound)
	res.Process()
}
