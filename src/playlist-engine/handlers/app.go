package handlers

import (
	"encoding/json"
	"net/http"
	"playlist-engine/models"
	"playlist-engine/spotify"
)

type App struct {
	Spotify *spotify.Client
}

func echoHeaders(w http.ResponseWriter, correlationID, experimentID string) {
	if correlationID != "" {
		w.Header().Set("X-Correlation-Id", correlationID)
	}
	if experimentID != "" {
		w.Header().Set("X-Experiment-Id", experimentID)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeProblem(w http.ResponseWriter, status int, title, detail, correlationID, experimentID string) {
	p := models.ProblemDetails{
		Type:          "https://tools.ietf.org/html/rfc9457",
		Title:         title,
		Status:        status,
		Detail:        detail,
		CorrelationID: correlationID,
		ExperimentID:  experimentID,
	}
	writeJSON(w, status, p)
}
