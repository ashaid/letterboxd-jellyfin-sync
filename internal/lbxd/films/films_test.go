package films

import (
	"net/http"
	"os"
	"testing"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/auth"
	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/utils"
	"github.com/joho/godotenv"
)

func TestUploadWatchlist(t *testing.T) {
	godotenv.Load("../../../.env")
	films, err := ReadFilmsFromCSV("../../../films_with_lids.csv")
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}

	lbxd := utils.NewLbxdClient(
		"https://api.letterboxd.com/api/v0",
		os.Getenv("LBXD_USER_NAME"),
		os.Getenv("LBXD_PASSWORD"),
		os.Getenv("LBXD_CLIENT_ID"),
		os.Getenv("LBXD_CLIENT_SECRET"),
	)

	tokens, err := auth.GetAccessTokens(lbxd.BaseURL, lbxd.Client_Id, lbxd.Username, lbxd.Password, lbxd.Client_Secret, lbxd.Client)
	if err != nil {
		t.Fatalf("Failed to get tokens: %v", err)
	}

	listResponse, statusCode, err := UploadAsWatchlist("https://letterboxd.com", &http.Client{}, tokens, films)
	if err != nil {
		t.Fatalf("Failed to upload: %v", err)
	}

	t.Logf("Status Code: %d", statusCode)
	t.Logf("Response: %+v", listResponse)
}
