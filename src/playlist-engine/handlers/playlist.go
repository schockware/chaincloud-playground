package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"playlist-engine/models"
	"playlist-engine/spotify"
	"strings"
	"time"
)

func (a *App) HandleGeneratePlaylist(w http.ResponseWriter, r *http.Request) {
	corrID := r.Header.Get("X-Correlation-Id")
	expID := r.Header.Get("X-Experiment-Id")
	echoHeaders(w, corrID, expID)

	var req models.GeneratePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeProblem(w, http.StatusBadRequest, "Invalid request body", err.Error(), corrID, expID)
		return
	}

	if !isValidCondition(req.Weather.Condition) || !isValidTimeOfDay(req.Weather.TimeOfDay) {
		writeProblem(w, http.StatusBadRequest, "Invalid weather fields",
			"condition and time_of_day must match allowed enum values", corrID, expID)
		return
	}

	experimentID := req.ExperimentID
	if experimentID == "" {
		experimentID = expID
	}

	mockReqHeader := r.Header.Get("X-ARBITRARY-MOCK")
	useMockSpotify := strings.Contains(mockReqHeader, "spotify")

	// In-process fallback: header says mock but no mock server is configured
	if useMockSpotify && a.SpotifyMock == nil {
		r8 := mapWeatherToRecipe(req.Weather)
		location := req.Weather.Description
		if req.LocationLabel != "" {
			location = req.LocationLabel
		}
		w.Header().Set("X-ARBITRARY-MOCK", "spotify")
		writeJSON(w, http.StatusOK, models.PlaylistResult{
			PlaylistID:      "mock-playlist-id",
			SpotifyEmbedURL: "https://open.spotify.com/embed/playlist/mock-playlist-id",
			PlaylistName:    capitalize(req.Weather.Condition) + " in " + location,
			Recipe:          r8,
			ExperimentID:    experimentID,
			TracksAdded:     5,
			CreatedAt:       time.Now().UTC(),
		})
		return
	}

	client := a.Spotify
	if useMockSpotify {
		client = a.SpotifyMock
	}

	ctx := r.Context()
	token, err := client.Token(ctx)
	if err != nil {
		writeProblem(w, http.StatusUnauthorized, "Spotify token unavailable", err.Error(), corrID, expID)
		return
	}

	userID, err := client.CurrentUserID(ctx, token)
	if err != nil {
		writeProblem(w, http.StatusBadGateway, "Spotify user lookup failed", err.Error(), corrID, expID)
		return
	}

	r8 := mapWeatherToRecipe(req.Weather)

	uris, err := client.SearchTracks(ctx, token, r8.Genres, r8.TrackCount)
	if err != nil {
		var rlErr *spotify.RateLimitError
		if errors.As(err, &rlErr) {
			w.Header().Set("Retry-After", rlErr.RetryAfter)
			writeProblem(w, http.StatusTooManyRequests, "Spotify rate limit exceeded", "", corrID, expID)
			return
		}
		writeProblem(w, http.StatusBadGateway, "Spotify search failed", err.Error(), corrID, expID)
		return
	}

	location := req.LocationLabel
	if location == "" {
		location = "Unknown Location"
	}
	playlistName := fmt.Sprintf("%s in %s", capitalize(req.Weather.Condition), location)
	description := fmt.Sprintf("Weather-driven playlist: %s, %s", req.Weather.Description, req.Weather.TimeOfDay)

	playlistID, err := client.CreatePlaylist(ctx, token, userID, playlistName, description)
	if err != nil {
		writeProblem(w, http.StatusBadGateway, "Spotify playlist creation failed", err.Error(), corrID, expID)
		return
	}

	if err := client.AddTracks(ctx, token, playlistID, uris); err != nil {
		writeProblem(w, http.StatusBadGateway, "Spotify add tracks failed", err.Error(), corrID, expID)
		return
	}

	if useMockSpotify {
		w.Header().Set("X-ARBITRARY-MOCK", "spotify")
	}

	writeJSON(w, http.StatusOK, models.PlaylistResult{
		PlaylistID:      playlistID,
		SpotifyEmbedURL: "https://open.spotify.com/embed/playlist/" + playlistID,
		PlaylistName:    playlistName,
		Recipe:          r8,
		ExperimentID:    experimentID,
		TracksAdded:     len(uris),
		CreatedAt:       time.Now().UTC(),
	})
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
