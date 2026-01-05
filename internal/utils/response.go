package utils

import (
	"encoding/json"
	"net/http"

	"github.com/fransiscushermanto/backend/internal/models"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		Log().Error().Err(err).Interface("payload", payload).Msg("Failed to marshal JSON response")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal server error"}`)) // Fallback generic error
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithSuccess(w http.ResponseWriter, code int, data interface{}, meta *map[string]interface{}) {
	result := models.ApiResult{
		Data: data,
		Meta: meta,
	}
	RespondWithJSON(w, code, result)
}
