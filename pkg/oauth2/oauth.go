package oauth2

import (
	"auth/internal/pkg/database"
	"auth/pkg/models"
	"auth/pkg/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func User(w http.ResponseWriter, r *http.Request) {
	cid := r.Header.Get("cid")
	user := models.User{}
	if err := database.DB.Where("id = ?", cid).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found. CID #%s.\n", cid)
			res := response.New(w, r, fmt.Sprintf("User not found. CID #%d.\n", cid), http.StatusNotFound)
			res.Process()
			return
		}

		log.Printf("Error occurred while executing the query on /api/user. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error while fetching user information.", http.StatusInternalServerError)
		res.Process()
		return
	}

	bytes, err := user.Json()

	if err != nil {
		log.Printf("Error occurred while marshalling the response on /api/user. Error: %s.", err.Error())
		res := response.New(w, r, "Internal server error while fetching user information.", http.StatusInternalServerError)
		res.Process()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	if _, err := w.Write(bytes); err != nil {
		log.Println(err.Error())
	}
}
