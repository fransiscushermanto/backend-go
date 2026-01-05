package controllers

import (
	"net/http"

	"github.com/fransiscushermanto/backend/internal/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
