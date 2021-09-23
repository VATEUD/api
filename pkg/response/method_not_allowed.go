package response

import (
	"log"
	"net/http"
)

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Page %s could got request with unsupported method - %s.\n", r.URL.String(), r.Method)
	res := New(w, r, "Method not allowed.", http.StatusMethodNotAllowed)
	res.Process()
}
