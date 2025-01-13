package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestArtits(t *testing.T) {
	emptyMock := func(mock sqlmock.Sqlmock) {}
	testCases := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		mockDbResponse func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "sort can be release_count",
			request:        createRequestForPath("/artists?artist=John Doe&sort=release_count&order=desc"),
			expectedStatus: http.StatusOK,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by release_count desc;")).
					WillReturnRows(standardArtistRows)
			},
		},
		{
			name:           "sort can be style_genre",
			request:        createRequestForPath("/artists?artist=John Doe&sort=style_genre&order=desc"),
			expectedStatus: http.StatusOK,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by style_genre desc;")).
					WillReturnRows(standardArtistRows)
			},
		},
		{
			name:           "sort can not be artist_name",
			request:        createRequestForPath("/artists?artist=John Doe&sort=artist_name&order=desc"),
			expectedStatus: http.StatusBadRequest,
			mockDbResponse: emptyMock,
		},
		{
			name:           "order can be asc",
			request:        createRequestForPath("/artists?artist=John Doe&sort=release_count&order=asc"),
			expectedStatus: http.StatusOK,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where artist_name = 'John Doe' order by release_count asc;")).
					WillReturnRows(standardArtistRows)
			},
		},
		{
			name:           "order can not be anything",
			request:        createRequestForPath("/artists?artist=John Doe&sort=release_count&order=anything"),
			expectedStatus: http.StatusBadRequest,
			mockDbResponse: emptyMock,
		},
		{
			name:           "when error occurs while fetching, 500 code is returned",
			request:        createRequestForPath("/artists?artist=John Doe&sort=release_count&order=asc"),
			expectedStatus: http.StatusInternalServerError,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from release_count rc where artist_name = ? order by ? asc")).
					WithArgs("John Doe", "release_count").
					WillReturnError(fmt.Errorf("test error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("unexpected error while opening mock db: %s", err)
			}
			defer db.Close()

			h := Handlers{
				Db: db,
			}

			tc.mockDbResponse(mock)

			rr := httptest.NewRecorder()
			testHandler := http.HandlerFunc(h.Artists)
			testHandler.ServeHTTP(rr, tc.request)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

		})
	}
}

func createRequestForPath(path string) *http.Request {
	req, _ := http.NewRequest("GET", path, nil)
	return req
}

var standardArtistRows = sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
	AddRow("John Doe", "Rock", 32).
	AddRow("John Doe", "Electronic", 23)
