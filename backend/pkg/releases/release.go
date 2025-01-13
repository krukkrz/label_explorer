package releases

type ReleaseCount struct {
	StyleGenre   string `json:"style_genre,omitempty"`
	ArtistName   string `json:"artist_name,omitempty"`
	ReleaseCount int    `json:"release_count,omitempty"`
}
