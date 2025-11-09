package main

import (
	"fmt"

	internal "github.com/ashaid/letterboxd-jellyfin-sync/internal/jellyfin"
)

func main() {
	// Get list of unwatched movies from jellyfin
	unwatchedMovies := internal.InvokeJellyfin()
	fmt.Print(unwatchedMovies)

	// Format into letterboxd csv
	// (Optional) rank movies
	// Upload to leeterboxd as watch list
}
