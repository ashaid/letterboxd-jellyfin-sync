package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/jellyfin/types"
	"github.com/joho/godotenv"
)

type JellyfinClient struct {
	BaseURL  string
	Username string
	Password string
	client   *http.Client
}

func (jc *JellyfinClient) getAccessToken() (string, string) {

	login := types.LoginRequest{
		Username: jc.Username,
		Pw:       jc.Password,
	}

	login_data, err := json.Marshal(login)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", jc.BaseURL+"/Users/AuthenticateByName", bytes.NewBuffer(login_data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", `MediaBrowser Client="letterboxd-jellyfin-sync", Device="Go", DeviceId="go-client-1", Version="1.0.0"`)

	res, err := jc.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		log.Fatalf("Status: %d, Body: %s", res.StatusCode, string(body))
	}

	var result map[string]any
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}

	accessToken := result["AccessToken"].(string)
	sessionInfo := result["SessionInfo"].(map[string]any)
	userId := sessionInfo["UserId"].(string)

	return accessToken, userId
}

func (jc *JellyfinClient) getMovies(accessToken string, userId string) {
	endpoint := fmt.Sprintf("%s/Users/%s/Items", jc.BaseURL, userId)

	params := url.Values{}
	params.Add("SortBy", "SortName,ProductionYear")
	params.Add("SortOrder", "Ascending")
	params.Add("IncludeItemTypes", "Movie")
	params.Add("Recursive", "true")

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	fmt.Print(fullURL)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(`MediaBrowser Client="letterboxd-jellyfin-sync", Device="Go", DeviceId="go-client-1", Version="1.0.0", Token="%s"`, accessToken))

	res, err := jc.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		log.Fatalf("Status: %d, Body: %s", res.StatusCode, string(body))
	}

	var result types.ItemsResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
}

func InvokeJellyfin() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	JELLYFIN_BASE_URL := os.Getenv("JELLYFIN_BASE_URL")
	JELLYFIN_USER_NAME := os.Getenv("JELLYFIN_USER_NAME")
	JELLYFIN_PASSWORD := os.Getenv("JELLYFIN_PASSWORD")

	client := &http.Client{}

	jc := &JellyfinClient{JELLYFIN_BASE_URL, JELLYFIN_USER_NAME, JELLYFIN_PASSWORD, client}

	accessToken, userId := jc.getAccessToken()

	jc.getMovies(accessToken, userId)

}
