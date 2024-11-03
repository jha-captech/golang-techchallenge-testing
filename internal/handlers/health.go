package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

// HandleHealthCheck handles the health check endpoint
//
//	@Summary		Health Check
//	@Description	Health Check endpoint
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	healthResponse
//	@Router			/health	[GET]
func HandleHealthCheck(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.InfoContext(r.Context(), "health check called")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "ok"})
	}
}
