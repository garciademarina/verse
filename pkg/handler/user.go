package handler

import (
	"log"
	"net/http"

	"github.com/garciademarina/verse/pkg/user"
)

// UserHandler handler struct for user endpoints
type UserHandler struct {
	repo user.Repository
}

// NewUserHandler create a new UserHandler
func NewUserHandler(repo user.Repository) UserHandler {
	return UserHandler{
		repo: repo,
	}
}

// FindById handles GET /xxx requests.
func (h *UserHandler) FindById(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] UserHandler.FindById\n")
		userID, err := GetJwtValue(r, "user_id")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "jwt_user_id_not_found"})
			return
		}

		payload, _ := h.repo.FindById(r.Context(), userID)

		respondwithJSON(w, http.StatusOK, payload)
	}
}
