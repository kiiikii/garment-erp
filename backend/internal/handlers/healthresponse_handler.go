package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/kiiikii/garment-erp/backend/internal/models"
)

func (app *App) HealthChecker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HealthResponse{
		Status: "success",
		Envi:   os.Getenv("APP_ENV"),
	}

	json.NewEncoder(w).Encode(response)
}
