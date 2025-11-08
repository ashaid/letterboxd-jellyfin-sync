package main

import (
	"fmt"

	internal "github.com/ashaid/letterboxd-jellyfin-sync/internal/jellyfin"
)

func main() {
	fmt.Println("Hello, World!")
	// Get list of unwatched movies from jellyfin
	internal.InvokeJellyfin()

	// create query to get movies by unwatched by user

	// Format into letterboxd csv
	// (Optional) rank movies
	// Upload to leeterboxd as watch list
}
