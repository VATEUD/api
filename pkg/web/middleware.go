package web

import (
	"api/internal/pkg/database"
	"api/pkg/jwt"
	"api/pkg/models"
	"api/pkg/response"
	"api/utils"
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const BaseRateLimit = 60

var limiter = rate.NewLimiter(rate.Every(time.Minute), rateLimit())

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.Split(r.RequestURI, "?")[0]

		if !server.NeedsAuth(uri) {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if server.GuestOnly(uri) {
			if authHeader != "" {
				log.Println("This route is for guests only.")
				res := response.New(w, r, "This route is for guests only.", http.StatusUnauthorized)
				res.Process()
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		if len(authHeader) < 1 {
			log.Println("Authentication header not provided.")
			res := response.New(w, r, "Authentication header not provided.", http.StatusUnauthorized)
			res.Process()
			return
		}

		if server.NeedsSubdivisionToken(uri) {
			token := strings.TrimPrefix(authHeader, "Token ")

			if len(token) < 1 {
				log.Println("Authentication header not provided.")
				res := response.New(w, r, "Authentication header not provided.", http.StatusUnauthorized)
				res.Process()
				return
			}

			var subToken models.SubdivisionToken

			if err := database.DB.Where("token = ?", token).First(&subToken).Error; err != nil {
				log.Println("Token not found.")
				res := response.New(w, r, "Authentication header not provided.", http.StatusUnauthorized)
				res.Process()
				return
			}

			ctx := context.WithValue(r.Context(), "token", subToken)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		auth := strings.TrimPrefix(authHeader, "Bearer ")

		if len(auth) < 1 {
			log.Println("Authentication header not provided.")
			res := response.New(w, r, "Authentication header not provided.", http.StatusUnauthorized)
			res.Process()
			return
		}

		token, err := jwt.New(auth)

		if err != nil {
			log.Println("Invalid token provided.")
			res := response.New(w, r, "Invalid token provided.", http.StatusUnauthorized)
			res.Process()
			return
		}

		if !token.Valid {
			log.Println("Invalid token provided.")
			res := response.New(w, r, "Invalid token provided.", http.StatusUnauthorized)
			res.Process()
			return
		}

		cid, ok := token.MapClaims["cid"]

		if !ok {
			log.Println("Invalid token provided. Token claims could not be parsed.")
			res := response.New(w, r, "Invalid token provided.", http.StatusUnauthorized)
			res.Process()
			return
		}

		r.Header.Set("cid", fmt.Sprintf("%v", cid))

		next.ServeHTTP(w, r)
	})
}

func rateLimitingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			log.Printf("Too many requests from the following IP %s.\n", r.Header.Get("IP"))
			res := response.New(w, r, "Too many requests.", http.StatusTooManyRequests)
			res.Process()
			return
		}

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.Split(r.RequestURI, "?")[0]

		if server.AllowCors(uri) {
			utils.Allow(w, "*")
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

	return BaseRateLimit
}
