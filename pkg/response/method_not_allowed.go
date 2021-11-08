package response

import (
	"api/internal/pkg/logger"
	"net/http"
)

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Printf("Page %s could got request with unsupported method - %s.\n", r.URL.String(), r.Method)
	res := New(w, r, "Method not allowed.", http.StatusMethodNotAllowed)
	res.Process()
}
