package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/ashaid/letterboxd-jellyfin-sync/types"
	"github.com/joho/godotenv"
)

func TestGetJellyfinUnwatchedMovies(t *testing.T) {
	t.Run("Env is properly configured", func(t *testing.T) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal("Error loading .env file")
		}
		JELLYFIN_BASE_URL := os.Getenv("JELLYFIN_BASE_URL")

		if JELLYFIN_BASE_URL == "" {
			t.Errorf("JELLYFIN_BASE_URL Not found in .env")
		}

		t.Log("env loaded")
	})

	t.Run("Check jellyfin server health", func(t *testing.T) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal("Error loading .env file")
		}
		JELLYFIN_BASE_URL := os.Getenv("JELLYFIN_BASE_URL")

		res, err := http.Get(JELLYFIN_BASE_URL + "/health")
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()

		if res.StatusCode > 299 {
			t.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s", body)
	})

	t.Run("Authenticate Jellyfin", func(t *testing.T) {
		err := godotenv.Load()
		if err != nil {
			t.Fatal("Error loading .env file")
		}
		JELLYFIN_BASE_URL := os.Getenv("JELLYFIN_BASE_URL")
		JELLYFIN_USER_NAME := os.Getenv("JELLYFIN_USER_NAME")
		JELLYFIN_PASSWORD := os.Getenv("JELLYFIN_PASSWORD")

		login := types.LoginRequest{
			Username: JELLYFIN_USER_NAME,
			Pw:       JELLYFIN_PASSWORD,
		}

		login_data, err := json.Marshal(login)
		if err != nil {
			t.Fatal(err)
		}
		client := &http.Client{}
		req, err := http.NewRequest("POST", JELLYFIN_BASE_URL+"/Users/AuthenticateByName", bytes.NewBuffer(login_data))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", `MediaBrowser Client="letterboxd-jellyfin-sync", Device="Go", DeviceId="go-client-1", Version="1.0.0"`)

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			body, _ := io.ReadAll(res.Body)
			t.Fatalf("Status: %d, Body: %s", res.StatusCode, string(body))
		}

		var result map[string]any
		json.NewDecoder(res.Body).Decode(&result)
		t.Logf("%s", result["AccessToken"])
	})

}
