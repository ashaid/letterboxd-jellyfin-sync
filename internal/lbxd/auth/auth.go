package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/types"
	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/utils"
)

func GetUnscopedToken(baseURL, clientID, clientSecret string, httpClient *http.Client) (*types.UnscopedTokenResponse, error) {
	endpoint := "/auth/token"
	requestURL := baseURL + endpoint

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("scope", "")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	utils.SetAPIHeaders(req)

	res, err := httpClient.Do(req)
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

func GetAccessTokens(baseURL, clientID, userName, password, clientSecret string, httpClient *http.Client) (*types.TokenResponse, error) {
	endpoint := "/auth/token"
	requestURL := baseURL + endpoint

	formData := url.Values{}
	formData.Set("grant_type", "password")
	formData.Set("username", userName)
	formData.Set("password", password)
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	utils.SetAPIHeaders(req)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("authentication failed with status %d: %s", res.StatusCode, string(body))
	}

	var tokenResponse types.TokenResponse
	if err := json.NewDecoder(res.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResponse, nil
}
