package response

import (
	"api/internal/pkg/logger"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Printf("Page %s could not be found.\n", r.URL.String())
	res := New(w, r, "The page you are looking for could not be found.", http.StatusNotFound)
	res.Process()
}
