package main

import (
	csv "github.com/ashaid/letterboxd-jellyfin-sync/internal/csv"
	jf "github.com/ashaid/letterboxd-jellyfin-sync/internal/jellyfin"
)

func main() {
	unwatchedMovies := jf.GetUnwatchedMoviesInCSVFormat()

	csv.CreateCSVInLetterboxdFormat(unwatchedMovies)
	// (Optional) rank movies
	// Upload to leeterboxd as watch list
}
