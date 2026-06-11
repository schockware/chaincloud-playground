package spotify

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type userProfile struct {
	ID string `json:"id"`
}

type searchResponse struct {
	Tracks struct {
		Items []trackItem `json:"items"`
	} `json:"tracks"`
}

type trackItem struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type createPlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type createdPlaylist struct {
	ID string `json:"id"`
}

type addTracksRequest struct {
	URIs []string `json:"uris"`
}
