package releases

import (
	"database/sql"
	"fmt"
)

func ByArtist(db *sql.DB, artistName string, sortBy string, order string) ([]ReleaseCount, error) {
	query := fmt.Sprintf("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = '%s' order by %s %s;", artistName, sortBy, order)
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query releases for %s: %v", artistName, err)
	}
	defer rows.Close()

	var results []ReleaseCount
	for rows.Next() {
		var rc ReleaseCount
		if err = rows.Scan(&rc.ArtistName, &rc.StyleGenre, &rc.ReleaseCount); err != nil {
			return nil, err
		}
		results = append(results, rc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
