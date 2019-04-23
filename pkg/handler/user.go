package handler

import (
	"log"
	"net/http"

	repository "github.com/garciademarina/verse/pkg/repository"
)

// UserHandler handler struct for user endpoints
type UserHandler struct {
	repo repository.UserRepo
}

// NewUserHandler create a new UserHandler
func NewUserHandler(repo repository.UserRepo) UserHandler {
	return UserHandler{
		repo: repo,
	}
}

// FindById handles GET /xxx requests.
func (h *UserHandler) FindById(logger *log.Logger) http.HandlerFunc {
	logger.Printf("[handle] UserHandler.FindById\n")
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetJwtValue(r, "user_id")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "jwt_user_id_not_found"})
			return
		}

		payload, _ := h.repo.FindById(r.Context(), userID)

		respondwithJSON(w, http.StatusOK, payload)
	}
}
