package models

import "time"

type WeatherCondition struct {
	Condition       string  `json:"condition"`
	TempCelsius     float64 `json:"temp_celsius"`
	FeelsLikeCelsius float64 `json:"feels_like_celsius"`
	Humidity        int     `json:"humidity"`
	WindSpeedMps    float64 `json:"wind_speed_mps"`
	Description     string  `json:"description"`
	TimeOfDay       string  `json:"time_of_day"`
}

type TempoRange struct {
	MinBPM int `json:"min_bpm"`
	MaxBPM int `json:"max_bpm"`
}

type PlaylistRecipe struct {
	Mood        string     `json:"mood"`
	Genres      []string   `json:"genres"`
	TempoRange  TempoRange `json:"tempo_range"`
	EnergyLevel float64    `json:"energy_level"`
	Valence     float64    `json:"valence"`
	TrackCount  int        `json:"track_count"`
}

type GeneratePlaylistRequest struct {
	Weather       WeatherCondition `json:"weather"`
	LocationLabel string           `json:"location_label,omitempty"`
	ExperimentID  string           `json:"experiment_id,omitempty"`
}

type PlaylistResult struct {
	PlaylistID      string         `json:"playlist_id"`
	SpotifyEmbedURL string         `json:"spotify_embed_url"`
	PlaylistName    string         `json:"playlist_name"`
	Recipe          PlaylistRecipe `json:"recipe"`
	ExperimentID    string         `json:"experiment_id,omitempty"`
	TracksAdded     int            `json:"tracks_added"`
	CreatedAt       time.Time      `json:"created_at"`
}

type HealthResponse struct {
	Status       string `json:"status"`
	SpotifyToken string `json:"spotify_token,omitempty"`
}

type ProblemDetails struct {
	Type          string `json:"type"`
	Title         string `json:"title"`
	Status        int    `json:"status"`
	Detail        string `json:"detail,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
	ExperimentID  string `json:"experiment_id,omitempty"`
}
