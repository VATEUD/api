package web

import (
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const BASE_RATE_LIMIT = 60

var limiter = rate.NewLimiter(rate.Every(time.Minute), rateLimit())

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.Split(r.RequestURI, "?")[0]

		if !server.NeedsAuth(uri) {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func rateLimit() int {
	if val := os.Getenv("RATE_LIMIT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}

	return BASE_RATE_LIMIT
}
