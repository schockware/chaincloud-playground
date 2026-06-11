package spotify

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	http           *http.Client
	clientID       string
	clientSecret   string
	refreshToken   string
	apiBase        string
	tokenURL       string
	mu             sync.Mutex
	accessToken    string
	tokenExpiresAt time.Time
}

func New(clientID, clientSecret, refreshToken, mockBaseURL, caFile string) *Client {
	apiBase := "https://api.spotify.com/v1"
	tokenURL := "https://accounts.spotify.com/api/token"
	if mockBaseURL != "" {
		apiBase = mockBaseURL + "/v1"
		tokenURL = mockBaseURL + "/api/token"
	}
	return &Client{
		http:         buildHTTPClient(caFile),
		clientID:     clientID,
		clientSecret: clientSecret,
		refreshToken: refreshToken,
		apiBase:      apiBase,
		tokenURL:     tokenURL,
	}
}

func buildHTTPClient(caFile string) *http.Client {
	if caFile == "" {
		return &http.Client{Timeout: 15 * time.Second}
	}
	pem, err := os.ReadFile(caFile)
	if err != nil {
		slog.Warn("could not read Spotify TLS CA file, using system pool", "file", caFile, "err", err)
		return &http.Client{Timeout: 15 * time.Second}
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(pem)
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		},
	}
}

func (c *Client) HasCredentials() bool {
	return c.clientID != "" && c.clientSecret != "" && c.refreshToken != ""
}

func (c *Client) Token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.tokenExpiresAt.Add(-30 * time.Second)) {
		return c.accessToken, nil
	}

	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {c.refreshToken},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	creds := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+creds)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("token refresh: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token refresh returned %d", resp.StatusCode)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("token decode: %w", err)
	}

	c.accessToken = tr.AccessToken
	c.tokenExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	if tr.RefreshToken != "" {
		c.refreshToken = tr.RefreshToken
	}
	return c.accessToken, nil
}

func (c *Client) CurrentUserID(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.apiBase+"/me", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get user returned %d", resp.StatusCode)
	}

	var profile userProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return "", fmt.Errorf("user decode: %w", err)
	}
	return profile.ID, nil
}

func (c *Client) SearchTracks(ctx context.Context, token string, genres []string, limit int) ([]string, error) {
	q := "genre:" + strings.Join(genres, " genre:")
	u := fmt.Sprintf("%s/search?q=%s&type=track&limit=%d&market=US",
		c.apiBase, url.QueryEscape(q), limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &RateLimitError{RetryAfter: resp.Header.Get("Retry-After")}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search returned %d", resp.StatusCode)
	}

	var sr searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("search decode: %w", err)
	}

	uris := make([]string, 0, len(sr.Tracks.Items))
	for _, item := range sr.Tracks.Items {
		uris = append(uris, item.URI)
	}
	return uris, nil
}

func (c *Client) CreatePlaylist(ctx context.Context, token, userID, name, description string) (string, error) {
	body, _ := json.Marshal(createPlaylistRequest{
		Name:        name,
		Description: description,
		Public:      false,
	})

	u := fmt.Sprintf("%s/users/%s/playlists", c.apiBase, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("create playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create playlist returned %d", resp.StatusCode)
	}

	var pl createdPlaylist
	if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
		return "", fmt.Errorf("create playlist decode: %w", err)
	}
	return pl.ID, nil
}

func (c *Client) AddTracks(ctx context.Context, token, playlistID string, uris []string) error {
	body, _ := json.Marshal(addTracksRequest{URIs: uris})

	u := fmt.Sprintf("%s/playlists/%s/tracks", c.apiBase, playlistID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("add tracks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("add tracks returned %d", resp.StatusCode)
	}
	return nil
}

type RateLimitError struct {
	RetryAfter string
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("spotify rate limit; retry after %s", e.RetryAfter)
}
