package handlers

import (
	"encoding/json"
	"net/http"
	"playlist-engine/models"
)

func (a *App) HandleRecipe(w http.ResponseWriter, r *http.Request) {
	corrID := r.Header.Get("X-Correlation-Id")
	expID := r.Header.Get("X-Experiment-Id")
	echoHeaders(w, corrID, expID)

	var weather models.WeatherCondition
	if err := json.NewDecoder(r.Body).Decode(&weather); err != nil {
		writeProblem(w, http.StatusBadRequest, "Invalid request body", err.Error(), corrID, expID)
		return
	}

	if !isValidCondition(weather.Condition) || !isValidTimeOfDay(weather.TimeOfDay) {
		writeProblem(w, http.StatusBadRequest, "Invalid weather fields",
			"condition and time_of_day must match allowed enum values", corrID, expID)
		return
	}

	writeJSON(w, http.StatusOK, mapWeatherToRecipe(weather))
}

func mapWeatherToRecipe(w models.WeatherCondition) models.PlaylistRecipe {
	switch w.Condition {
	case "clear":
		switch w.TimeOfDay {
		case "morning":
			return recipe("upbeat", []string{"pop", "indie pop"}, 110, 140, 0.80, 0.80, 20)
		case "afternoon":
			return recipe("happy", []string{"pop", "dance"}, 120, 150, 0.85, 0.90, 20)
		case "evening":
			return recipe("calm", []string{"indie", "acoustic"}, 80, 110, 0.50, 0.70, 15)
		default:
			return recipe("dark", []string{"electronic", "ambient"}, 70, 100, 0.40, 0.30, 15)
		}
	case "clouds":
		return recipe("calm", []string{"indie", "folk"}, 80, 110, 0.50, 0.50, 15)
	case "rain", "drizzle":
		switch w.TimeOfDay {
		case "morning":
			return recipe("melancholic", []string{"jazz", "acoustic"}, 70, 100, 0.40, 0.30, 15)
		case "afternoon":
			return recipe("cozy", []string{"indie pop", "folk"}, 80, 110, 0.50, 0.50, 15)
		case "evening":
			return recipe("melancholic", []string{"piano", "classical"}, 60, 90, 0.30, 0.30, 15)
		default:
			return recipe("dark", []string{"ambient", "electronic"}, 60, 90, 0.35, 0.25, 10)
		}
	case "thunderstorm":
		return recipe("dramatic", []string{"metal", "rock"}, 130, 170, 0.90, 0.20, 20)
	case "snow":
		return recipe("cozy", []string{"classical", "jazz"}, 70, 100, 0.40, 0.50, 15)
	default: // mist, fog, haze
		return recipe("calm", []string{"ambient", "electronic"}, 70, 100, 0.30, 0.40, 10)
	}
}

func recipe(mood string, genres []string, minBPM, maxBPM int, energy, valence float64, trackCount int) models.PlaylistRecipe {
	return models.PlaylistRecipe{
		Mood:        mood,
		Genres:      genres,
		TempoRange:  models.TempoRange{MinBPM: minBPM, MaxBPM: maxBPM},
		EnergyLevel: energy,
		Valence:     valence,
		TrackCount:  trackCount,
	}
}

func isValidCondition(c string) bool {
	switch c {
	case "clear", "clouds", "rain", "drizzle", "thunderstorm", "snow", "mist", "fog", "haze":
		return true
	}
	return false
}

func isValidTimeOfDay(t string) bool {
	switch t {
	case "morning", "afternoon", "evening", "night":
		return true
	}
	return false
}
