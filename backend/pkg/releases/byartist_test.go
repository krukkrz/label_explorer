package releases

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/stretchr/testify/assert"
)

func TestByArtist(t *testing.T) {
	type params struct {
		artistName string
		sortBy     string
		order      string
	}

	testCases := []struct {
		name string
		params
		query          string
		rows           *sqlmock.Rows
		expected       []ReleaseCount
		expectingError bool
	}{
		{
			name: "sorts by release count desc",
			params: params{
				artistName: "John Doe",
				sortBy:     "release_count",
				order:      "desc",
			},
			query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by release_count desc;",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("John Doe", "Rock", 32).
				AddRow("John Doe", "Electronic", 23),
			expected: []ReleaseCount{
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Rock",
					ReleaseCount: 32,
				},
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Electronic",
					ReleaseCount: 23,
				},
			},
		},
		{
			name: "sorts by release count asc",
			params: params{
				artistName: "John Doe",
				sortBy:     "release_count",
				order:      "asc",
			}, query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by release_count asc",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("John Doe", "Electronic", 23).
				AddRow("John Doe", "Rock", 32),
			expected: []ReleaseCount{
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Electronic",
					ReleaseCount: 23,
				},
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Rock",
					ReleaseCount: 32,
				},
			},
		},
		{
			name: "sorts by style/genre desc",
			params: params{
				artistName: "John Doe",
				sortBy:     "style_genre",
				order:      "desc",
			}, query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by style_genre desc",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("John Doe", "Rock", 32).
				AddRow("John Doe", "Electronic", 23),
			expected: []ReleaseCount{
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Rock",
					ReleaseCount: 32,
				},
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Electronic",
					ReleaseCount: 23,
				},
			},
		},
		{
			name: "sorts by style/genre asc",
			params: params{
				artistName: "John Doe",
				sortBy:     "style_genre",
				order:      "asc",
			}, query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by style_genre asc",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("John Doe", "Electronic", 23).
				AddRow("John Doe", "Rock", 32),
			expected: []ReleaseCount{
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Electronic",
					ReleaseCount: 23,
				},
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Rock",
					ReleaseCount: 32,
				},
			},
		},
		{
			name: "when query returns an error it returns it",
			params: params{
				artistName: "John Doe",
				sortBy:     "style_genre",
				order:      "asc",
			},
			query:          "select * from release_count where artist_name = ? order by ? asc",
			expectingError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unexpected error while opening mock db: %s", err)
			}
			defer db.Close()

			if tc.expectingError {
				mock.ExpectQuery(regexp.QuoteMeta(tc.query)).
					WillReturnError(fmt.Errorf("test error"))

			} else {
				mock.ExpectQuery(regexp.QuoteMeta(tc.query)).
					WillReturnRows(tc.rows)
			}

			actual, err := ByArtist(db, tc.artistName, tc.sortBy, tc.order)
			if tc.expectingError {
				if err == nil {
					t.Errorf("expecting error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				assert.Equal(t, tc.expected, actual)
			}

		})
	}
}
