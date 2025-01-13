package discogs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleaseDetails(t *testing.T) {
	testCases := []struct {
		name          string
		handler       http.HandlerFunc
		expected      DetailsResponse
		errorExpected bool
	}{
		{
			name: "fetches release details from discogs API",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(DetailsResponse{
					Id: 1,
					Genres: []string{
						"Electronic",
					},
					Styles: []string{
						"Downtempo",
						"Ethereal",
						"Contemporary R&B",
					},
				})
			},
			expected: DetailsResponse{
				Id: 1,
				Genres: []string{
					"Electronic",
				},
				Styles: []string{
					"Downtempo",
					"Ethereal",
					"Contemporary R&B",
				},
			},
		},
		{
			name: "if API response code is other than 200, the error is returned",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(DetailsResponse{})
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			httpClient := &http.Client{}
			client := &Client{
				HttpClient: httpClient,
			}

			actual, err := ReleaseDetails(client, server.URL)

			if tc.errorExpected {
				if err == nil {
					t.Errorf("error is expected but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("error ocurred but not expected: %v", err)
			}

			assert.Equal(t, actual, tc.expected)
		})
	}
}
