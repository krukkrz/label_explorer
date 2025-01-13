package releases

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestByStyleGenre(t *testing.T) {
	type params struct {
		styleGenreName string
		sortBy         string
		order          string
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
				styleGenreName: "Electronic",
				sortBy:         "release_count",
				order:          "desc",
			},
			query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by release_count desc;",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("Popeye", "Electronic", 32).
				AddRow("John Doe", "Electronic", 23),
			expected: []ReleaseCount{
				{
					ArtistName:   "Popeye",
					StyleGenre:   "Electronic",
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
				styleGenreName: "Electronic",
				sortBy:         "release_count",
				order:          "asc",
			},
			query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by release_count asc;",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("John Doe", "Electronic", 23).
				AddRow("Popeye", "Electronic", 32),
			expected: []ReleaseCount{
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Electronic",
					ReleaseCount: 23,
				},
				{
					ArtistName:   "Popeye",
					StyleGenre:   "Electronic",
					ReleaseCount: 32,
				},
			},
		},
		{
			name: "sorts by artist name desc",
			params: params{
				styleGenreName: "Electronic",
				sortBy:         "artist_name",
				order:          "desc",
			},
			query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by artist_name desc;",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("Popeye", "Electronic", 32).
				AddRow("John Doe", "Electronic", 23),
			expected: []ReleaseCount{
				{
					ArtistName:   "Popeye",
					StyleGenre:   "Electronic",
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
			name: "sorts by artist name asc",
			params: params{
				styleGenreName: "Electronic",
				sortBy:         "artist_name",
				order:          "asc",
			},
			query: "select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by artist_name asc;",
			rows: sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
				AddRow("John Doe", "Electronic", 23).
				AddRow("Popeye", "Electronic", 32),
			expected: []ReleaseCount{
				{
					ArtistName:   "John Doe",
					StyleGenre:   "Electronic",
					ReleaseCount: 23,
				},
				{
					ArtistName:   "Popeye",
					StyleGenre:   "Electronic",
					ReleaseCount: 32,
				},
			},
		},
		{
			name: "when query returns an error it returns it",
			params: params{
				styleGenreName: "Electronic",
				sortBy:         "release_count",
				order:          "asc",
			}, query: "select * from release_count where style_genre = ? order by ? asc",
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

			actual, err := ByStyleGenre(db, tc.styleGenreName, tc.sortBy, tc.order)
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
