package main

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestEnvAndJellyfinConnection(t *testing.T) {
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

}
