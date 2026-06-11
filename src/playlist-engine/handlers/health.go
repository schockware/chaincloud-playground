package handlers

import (
	"net/http"
	"playlist-engine/models"
)

func (a *App) HandleHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tokenStatus := "ok"

	if _, err := a.Spotify.Token(ctx); err != nil {
		if a.Spotify.HasCredentials() {
			tokenStatus = "expired"
		} else {
			tokenStatus = "missing"
		}
	}

	status := "ok"
	if tokenStatus != "ok" {
		status = "degraded"
	}

	writeJSON(w, http.StatusOK, models.HealthResponse{
		Status:       status,
		SpotifyToken: tokenStatus,
	})
}
