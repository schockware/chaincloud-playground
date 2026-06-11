package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"playlist-engine/handlers"
	"playlist-engine/spotify"
)

func main() {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		slog.Warn("Spotify credentials not fully configured; /playlist/generate will return 401")
	}

	app := &handlers.App{
		Spotify: spotify.New(clientID, clientSecret, refreshToken),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /playlist/generate", app.HandleGeneratePlaylist)
	mux.HandleFunc("POST /recipe", app.HandleRecipe)
	mux.HandleFunc("GET /health", app.HandleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5100"
	}

	addr := ":" + port
	slog.Info("playlist-engine listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
