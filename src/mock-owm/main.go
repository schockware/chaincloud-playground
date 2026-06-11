package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type owmWeather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

type owmMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
}

type owmWind struct {
	Speed float64 `json:"speed"`
}

type owmResponse struct {
	Weather  []owmWeather `json:"weather"`
	Main     owmMain      `json:"main"`
	Wind     owmWind      `json:"wind"`
	Timezone int64        `json:"timezone"`
	Dt       int64        `json:"dt"`
}

var conditions = [7]owmResponse{
	{Weather: []owmWeather{{"Clear", "clear sky"}}, Main: owmMain{22.5, 21.0, 55}, Wind: owmWind{3.2}},
	{Weather: []owmWeather{{"Clouds", "overcast clouds"}}, Main: owmMain{18.0, 17.0, 70}, Wind: owmWind{4.5}},
	{Weather: []owmWeather{{"Rain", "moderate rain"}}, Main: owmMain{14.0, 13.0, 85}, Wind: owmWind{5.0}},
	{Weather: []owmWeather{{"Thunderstorm", "thunderstorm with rain"}}, Main: owmMain{16.0, 15.0, 88}, Wind: owmWind{9.0}},
	{Weather: []owmWeather{{"Snow", "light snow"}}, Main: owmMain{-2.0, -5.0, 80}, Wind: owmWind{3.0}},
	{Weather: []owmWeather{{"Mist", "mist"}}, Main: owmMain{12.0, 11.5, 92}, Wind: owmWind{1.5}},
	{Weather: []owmWeather{{"Drizzle", "light intensity drizzle"}}, Main: owmMain{13.0, 12.5, 90}, Wind: owmWind{2.5}},
}

func main() {
	certFile := envOr("MOCK_TLS_CERT_FILE", "certs/mock-owm.crt")
	keyFile := envOr("MOCK_TLS_KEY_FILE", "certs/mock-owm.key")
	port := envOr("MOCK_PORT", "5300")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /data/2.5/weather", handleWeather)

	addr := ":" + port
	slog.Info("mock-owm listening", "addr", addr)
	if err := http.ListenAndServeTLS(addr, certFile, keyFile, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	hour := time.Now().UTC().Hour()

	raw := int(lat*10) + int(lon*10) + hour
	if raw < 0 {
		raw = -raw
	}
	idx := raw % len(conditions)

	resp := conditions[idx]
	resp.Dt = time.Now().Unix()
	resp.Timezone = 0

	writeJSON(w, http.StatusOK, resp)
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
