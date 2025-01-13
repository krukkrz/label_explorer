package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStyles(t *testing.T) {
	emptyMock := func(mock sqlmock.Sqlmock) {}
	testCases := []struct {
		name           string
		request        *http.Request
		expectedStatus int
		mockDbResponse func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "sort can be release_count",
			request:        createRequestForPath("/styles?style=Electronic&sort=release_count&order=desc"),
			expectedStatus: http.StatusOK,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by release_count desc;")).
					WillReturnRows(standardStyleRows)
			},
		},
		{
			name:           "sort can be artist_name",
			request:        createRequestForPath("/styles?style=Electronic&sort=artist_name&order=desc"),
			expectedStatus: http.StatusOK,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by artist_name desc;")).
					WillReturnRows(standardStyleRows)
			},
		},
		{
			name:           "sort can not be style_genre",
			request:        createRequestForPath("/styles?style=Electronic&sort=style_genre&order=desc"),
			expectedStatus: http.StatusBadRequest,
			mockDbResponse: emptyMock,
		},
		{
			name:           "order can be asc",
			request:        createRequestForPath("/styles?style=Electronic&sort=release_count&order=asc"),
			expectedStatus: http.StatusOK,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from releases_count rc where style_genre = 'Electronic' order by release_count asc")).
					WillReturnRows(standardArtistRows)
			},
		},
		{
			name:           "order can not be anything",
			request:        createRequestForPath("/styles?style=Electronic&sort=release_count&order=anything"),
			expectedStatus: http.StatusBadRequest,
			mockDbResponse: emptyMock,
		},
		{
			name:           "when error occurs while fetching, 500 code is returned",
			request:        createRequestForPath("/styles?style=Electronic&sort=release_count&order=asc"),
			expectedStatus: http.StatusInternalServerError,
			mockDbResponse: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("select rc.artist_name, rc.style_genre, rc.release_count from release_count rc where style_genre = ? order by ? asc")).
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
			testHandler := http.HandlerFunc(h.Styles)
			testHandler.ServeHTTP(rr, tc.request)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

		})
	}
}

var standardStyleRows = sqlmock.NewRows([]string{"artist_name", "style_genre", "release_count"}).
	AddRow("John Doe", "Electronic", 32).
	AddRow("Popeye", "Electronic", 23)
