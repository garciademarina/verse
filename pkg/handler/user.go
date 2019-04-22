package handler

import (
	"log"
	"net/http"

	repository "github.com/garciademarina/verse/pkg/repository"
)

// User ...
type User struct {
	repo repository.UserRepo
}

// NewUserHandler explain ...
func NewUserHandler(repo repository.UserRepo) User {
	return User{
		repo: repo,
	}
}

// FindById handles GET /top requests.
func (u *User) FindById(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := GetJwtValue(r, "user_id")

		err := r.ParseForm()
		if err != nil {
			logger.Printf("Error parsing the request form %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userID == "" {
			logger.Printf("userID parameter not found in Jwt\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logger.Printf("...handling FindById %s\n", userID)
		payload, _ := u.repo.FindById(r.Context(), userID)

		respondwithJSON(w, http.StatusOK, payload)
	}
}
