package utils

import (
	"fmt"
	"strings"
)

func BuildWatchlistPayload(lids []string) string {
	entries := make([]string, 0, len(lids))
	for _, lid := range lids {
		entries = append(entries, fmt.Sprintf(`{"film":"%s","action":"ADD","containsSpoilers":false}`, lid))
	}
	return fmt.Sprintf(
		`{"version":"0","name":"Jellyfin+Watchlist","description":"Synced+from+Jellyfin","tags":[],"published":false,"sharePolicy":"You","ranked":false,"entries":[%s]}`,
		strings.Join(entries, ","),
	)
}
