package letterboxd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/letterboxd/types"
)

type LetterboxdClient struct {
	BaseURL       string
	Username      string
	Password      string
	Client_Id     string
	Client_Secret string
	client        *http.Client
}

func newClient(baseURL, username, password, clientID, clientSecret string) *LetterboxdClient {
	return &LetterboxdClient{
		BaseURL:       baseURL,
		Username:      username,
		Password:      password,
		Client_Id:     clientID,
		Client_Secret: clientSecret,
		client:        &http.Client{},
	}
}

func (lbxd *LetterboxdClient) setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Letterboxd/6553 CFNetwork/3826.600.41 Darwin/24.6.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
}

func (lbxd *LetterboxdClient) getUnscopedToken() (*types.UnscopedTokenResponse, error) {
	endpoint := "/auth/token"
	requestURL := lbxd.BaseURL + endpoint

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("scope", "")
	formData.Set("client_id", lbxd.Client_Id)
	formData.Set("client_secret", lbxd.Client_Secret)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	lbxd.setCommonHeaders(req)

	res, err := lbxd.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("authentication failed with status %d: %s", res.StatusCode, string(body))
	}

	var tokenResponse types.UnscopedTokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResponse, nil
}

func Main() {
	lbxd := newClient(
		"https://api.letterboxd.com/api/v0",
		os.Getenv("LBXD_USER_NAME"),
		os.Getenv("LBXD_PASSWORD"),
		os.Getenv("LBXD_CLIENT_ID"),
		os.Getenv("LBXD_CLIENT_SECRET"),
	)

	_, err := lbxd.getUnscopedToken()
	if err != nil {
		log.Fatalf("Failed to retreive unscoped token: %v", err)
	}

	fmt.Printf("Successfully retrieved unscoped token")

	// lbxd.uploadAsWatchList()
}
