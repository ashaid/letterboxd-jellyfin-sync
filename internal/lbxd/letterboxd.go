package letterboxd

import (
	"fmt"
	"log"
	"os"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/auth"
	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/films"
	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/utils"
)

func Main() {
	lbxd := utils.NewLbxdClient(
		"https://api.letterboxd.com/api/v0",
		os.Getenv("LBXD_USER_NAME"),
		os.Getenv("LBXD_PASSWORD"),
		os.Getenv("LBXD_CLIENT_ID"),
		os.Getenv("LBXD_CLIENT_SECRET"),
	)

	_, err := auth.GetUnscopedToken(lbxd.BaseURL, lbxd.Client_Id, lbxd.Client_Secret, lbxd.Client)
	if err != nil {
		log.Fatalf("Failed to retrieve unscoped token: %v", err)
	}
	fmt.Printf("Successfully retrieved unscoped token\n")

	tokens, err := auth.GetAccessTokens(lbxd.BaseURL, lbxd.Client_Id, lbxd.Username, lbxd.Password, lbxd.Client_Secret, lbxd.Client)
	if err != nil {
		log.Fatalf("Failed to retrieve access tokens: %v", err)
	}
	fmt.Printf("Successfully retrieved access tokens\n")

	simpleClient := utils.NewSimpleClient("https://letterboxd.com")

	result, err := films.ProcessFilms("result.csv", simpleClient)
	if err != nil {
		log.Fatalf("Failed to process films: %v", err)
	}

	fmt.Printf("\n=== Processing Summary ===\n")
	fmt.Printf("Total films processed: %d\n", result.TotalProcessed)
	fmt.Printf("Successful: %d\n", len(result.SuccessfulFilms))
	fmt.Printf("Failed: %d\n", len(result.FailedFilms))

	films.UploadAsWatchlist(lbxd.BaseURL, lbxd.Client, tokens, result.SuccessfulFilms)
}
