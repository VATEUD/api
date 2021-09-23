package utils

import "net/http"

func Allow(w http.ResponseWriter, sources ...string) {
	for _, source := range sources {
		w.Header().Set("Access-Control-Allow-Origin", source)
	}
}
