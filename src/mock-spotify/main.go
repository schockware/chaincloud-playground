package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	certFile := envOr("MOCK_TLS_CERT_FILE", "certs/mock-spotify.crt")
	keyFile := envOr("MOCK_TLS_KEY_FILE", "certs/mock-spotify.key")
	port := envOr("MOCK_PORT", "5200")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/token", handleToken)
	mux.HandleFunc("GET /v1/me", handleMe)
	mux.HandleFunc("GET /v1/search", handleSearch)
	mux.HandleFunc("POST /v1/users/{user_id}/playlists", handleCreatePlaylist)
	mux.HandleFunc("POST /v1/playlists/{playlist_id}/tracks", handleAddTracks)

	addr := ":" + port
	slog.Info("mock-spotify listening", "addr", addr)
	if err := http.ListenAndServeTLS(addr, certFile, keyFile, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func handleToken(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": "mock-token",
		"token_type":   "Bearer",
		"expires_in":   3600,
	})
}

func handleMe(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"id":           "mock-user",
		"display_name": "Mock User",
	})
}

func handleSearch(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"tracks": map[string]any{
			"items": []map[string]any{
				{"uri": "spotify:track:mock001", "name": "Mock Track 1", "artists": []map[string]any{{"name": "Mock Artist"}}},
				{"uri": "spotify:track:mock002", "name": "Mock Track 2", "artists": []map[string]any{{"name": "Mock Artist"}}},
				{"uri": "spotify:track:mock003", "name": "Mock Track 3", "artists": []map[string]any{{"name": "Mock Artist"}}},
				{"uri": "spotify:track:mock004", "name": "Mock Track 4", "artists": []map[string]any{{"name": "Mock Artist"}}},
				{"uri": "spotify:track:mock005", "name": "Mock Track 5", "artists": []map[string]any{{"name": "Mock Artist"}}},
			},
		},
	})
}

func handleCreatePlaylist(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusCreated, map[string]any{
		"id": "mock-playlist-id",
		"external_urls": map[string]string{
			"spotify": "https://open.spotify.com/playlist/mock",
		},
	})
}

func handleAddTracks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusCreated, map[string]any{
		"snapshot_id": "mock-snapshot",
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
