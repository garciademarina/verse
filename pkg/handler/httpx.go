package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
)

// respondwithJSON write json response format
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// RespondWithError replies to the request with the specified error message (as a json) and HTTP code.
func RespondWithError(w http.ResponseWriter, code int, payload APIError) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// GetJwtValue returns jwt value for a given key
func GetJwtValue(r *http.Request, key string) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", err
	}
	keyValue, ok := claims[key].(string)
	if !ok {
		return "", errors.New("jwt key not found")
	}

	return keyValue, nil
}

// APIError represents api error messages
type APIError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
