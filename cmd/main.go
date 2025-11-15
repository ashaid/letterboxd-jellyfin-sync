package main

import (
	csv "github.com/ashaid/letterboxd-jellyfin-sync/internal/csv"
	jf "github.com/ashaid/letterboxd-jellyfin-sync/internal/jellyfin"
	lbxd "github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	unwatchedMovies := jf.GetUnwatchedMoviesInCSVFormat()

	csv.CreateCSVInLetterboxdFormat(unwatchedMovies)
	// (Optional) rank movies
	// Upload to leeterboxd as watch list

	lbxd.Main()
}
