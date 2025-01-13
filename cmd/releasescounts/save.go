package releasescounts

import (
	"database/sql"
	"fmt"
)

func Save(
	db *sql.DB,
	releasesCounts []ReleaseCount,
) error {
	query := `
	INSERT INTO releases_count (label_id, style_genre, artist_name, release_count)
	VALUES ($1, $2, $3, $4);
	`

	for _, releaseCount := range releasesCounts {
		_, err := db.Exec(query, releaseCount.LabelId, releaseCount.StyleGenre, releaseCount.ArtistName, releaseCount.ReleaseCount)
		if err != nil {
			return fmt.Errorf("failed to store release count: %w", err)
		}
	}

	return nil
}

type ReleaseCount struct {
	LabelId      int
	ArtistName   string
	StyleGenre   string
	ReleaseCount int
}
