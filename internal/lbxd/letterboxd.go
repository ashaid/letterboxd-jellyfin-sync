package letterboxd

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/auth"
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

	_, err = auth.GetAccessTokens(lbxd.BaseURL, lbxd.Client_Id, lbxd.Username, lbxd.Password, lbxd.Client_Secret, lbxd.Client)
	if err != nil {
		log.Fatalf("Failed to retrieve access tokens: %v", err)
	}

	fmt.Printf("Successfully retrieved access tokens\n")

	simpleClient := utils.NewSimpleClient(
		"https://letterboxd.com",
	)

	body, err := simpleClient.Get("/film/" + "10-things-i-hate-about-you")
	if err != nil {
		log.Fatalf("Failed to retrieve film: %v", err)
	}

	lidPattern := regexp.MustCompile(`"lid":\s*"([^"]+)"`)
	matches := lidPattern.FindSubmatch(body)
	if len(matches) < 2 {
		log.Fatalf("Failed to find lid in response")
	}

	lid := string(matches[1])
	fmt.Printf("Extracted LID: %s\n", lid)

	// lbxd.uploadAsWatchList()
}
