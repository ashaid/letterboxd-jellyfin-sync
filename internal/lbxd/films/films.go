package films

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/ashaid/letterboxd-jellyfin-sync/internal/lbxd/utils"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Film struct {
	Title string
	Year  string
	Slug  string
	LID   string
}

func TitleToSlug(title string) string {
	slug := strings.ToLower(title)

	// Handle special character ½
	slug = strings.ReplaceAll(slug, "½", "-half")

	// Remove accents from characters (é→e, á→a, ï→i, etc.)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	slug, _, _ = transform.String(t, slug)

	// Remove apostrophes
	slug = strings.ReplaceAll(slug, "'", "")

	// Remove colons
	slug = strings.ReplaceAll(slug, ":", "")

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

func ReadFilmsFromCSV(filepath string) ([]Film, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file is empty or has no data rows")
	}

	var films []Film
	for i, record := range records {
		// Skip header row
		if i == 0 {
			continue
		}

		if len(record) < 2 {
			continue
		}

		title := record[0]
		year := record[1]
		slug := TitleToSlug(title)

		films = append(films, Film{
			Title: title,
			Year:  year,
			Slug:  slug,
		})
	}

	return films, nil
}

func WriteFilmsToCSV(filepath string, films []Film) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Title", "Year", "Slug", "LID"}); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, film := range films {
		if err := writer.Write([]string{film.Title, film.Year, film.Slug, film.LID}); err != nil {
			return fmt.Errorf("failed to write film record: %w", err)
		}
	}

	return nil
}

func GetFilmId(simpleClient *utils.SimpleClient, filmSlug string) (string, error) {
	body, err := simpleClient.Get("/film/" + filmSlug)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve film: %w", err)
	}

	lidPattern := regexp.MustCompile(`"lid":\s*"([^"]+)"`)
	matches := lidPattern.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to find lid in response")
	}

	lid := string(matches[1])
	return lid, nil
}

func GetFilmIdWithYear(simpleClient *utils.SimpleClient, filmSlug, year string) (string, error) {
	// Try without year first
	lid, err := GetFilmId(simpleClient, filmSlug)
	if err == nil {
		return lid, nil
	}

	// Try with year appended
	slugWithYear := filmSlug + "-" + year
	lid, err = GetFilmId(simpleClient, slugWithYear)
	if err != nil {
		return "", fmt.Errorf("failed with both slug variations: %w", err)
	}

	return lid, nil
}

type ProcessingResult struct {
	SuccessfulFilms []Film
	FailedFilms     []Film
	TotalProcessed  int
}

func ProcessFilms(inputCSV string, simpleClient *utils.SimpleClient) (*ProcessingResult, error) {
	filmsList, err := ReadFilmsFromCSV(inputCSV)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	fmt.Printf("Read %d films from CSV\n", len(filmsList))

	err = WriteFilmsToCSV("films_with_slugs.csv", filmsList)
	if err != nil {
		return nil, fmt.Errorf("failed to write films with slugs: %w", err)
	}

	fmt.Printf("Saved films with slugs to films_with_slugs.csv\n")

	var successfulFilms []Film
	var failedFilms []Film

	for i := range filmsList {
		lid, err := GetFilmIdWithYear(simpleClient, filmsList[i].Slug, filmsList[i].Year)
		if err != nil {
			fmt.Printf("Failed to get film ID for '%s' (%s): %v\n", filmsList[i].Title, filmsList[i].Slug, err)
			failedFilms = append(failedFilms, filmsList[i])
			continue
		}

		filmsList[i].LID = lid
		successfulFilms = append(successfulFilms, filmsList[i])
		fmt.Printf("Film: %s (%s) -> Slug: %s -> LID: %s\n", filmsList[i].Title, filmsList[i].Year, filmsList[i].Slug, lid)
	}

	if len(successfulFilms) > 0 {
		err = WriteFilmsToCSV("films_with_lids.csv", successfulFilms)
		if err != nil {
			return nil, fmt.Errorf("failed to write successful films: %w", err)
		}
		fmt.Printf("Saved %d films with LIDs to films_with_lids.csv\n", len(successfulFilms))
	}

	if len(failedFilms) > 0 {
		err = WriteFilmsToCSV("films_failed.csv", failedFilms)
		if err != nil {
			return nil, fmt.Errorf("failed to write failed films: %w", err)
		}
		fmt.Printf("Saved %d failed films to films_failed.csv\n", len(failedFilms))
	}

	return &ProcessingResult{
		SuccessfulFilms: successfulFilms,
		FailedFilms:     failedFilms,
		TotalProcessed:  len(filmsList),
	}, nil
}
