package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func CreateCSVInLetterboxdFormat(unwatchedMovies [][]string) {
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	writer.WriteAll(unwatchedMovies)
}
