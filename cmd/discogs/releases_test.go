package discogs

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReleases(t *testing.T) {
	testCases := []struct {
		name          string
		handler       func() http.HandlerFunc
		expected      []Release
		errorExpected bool
	}{
		{
			name: "if there there is a single page, it calls API just once",
			handler: func() http.HandlerFunc {
				responses := []ReleasesResponse{
					{
						Pagination: Pagination{
							Page:    1,
							Pages:   1,
							PerPage: 2,
							Items:   2,
						},
						Releases: []Release{
							{
								Id:          1,
								Artist:      "John Doe",
								ResourceUrl: "http://somehting.com/release/1",
							},
							{
								Id:          2,
								Artist:      "Marry Poppins",
								ResourceUrl: "http://somehting.com/release/2",
							},
						},
					},
				}
				callCount := 0
				handlerFunc := func(w http.ResponseWriter, _ *http.Request) {
					if callCount >= len(responses) {
						t.Fatalf("unexpected call to handler: %d", callCount)
					}

					resp := responses[callCount]
					callCount++
					w.WriteHeader(200)
					json.NewEncoder(w).Encode(resp)
				}
				return handlerFunc
			},
			expected: []Release{
				{
					Id:          1,
					Artist:      "John Doe",
					ResourceUrl: "http://somehting.com/release/1",
				},
				{
					Id:          2,
					Artist:      "Marry Poppins",
					ResourceUrl: "http://somehting.com/release/2",
				},
			},
		},
		{
			name: "if there are more pages than one it fetch them all",
			handler: func() http.HandlerFunc {
				responses := []ReleasesResponse{
					{
						Pagination: Pagination{
							Page:    1,
							Pages:   2,
							PerPage: 2,
							Items:   3,
						},
						Releases: []Release{
							{
								Id:          1,
								Artist:      "John Doe",
								ResourceUrl: "http://somehting.com/release/1",
							},
							{
								Id:          2,
								Artist:      "Marry Poppins",
								ResourceUrl: "http://somehting.com/release/2",
							},
						},
					},
					{
						Pagination: Pagination{
							Page:    2,
							Pages:   2,
							PerPage: 2,
							Items:   3,
						},
						Releases: []Release{
							{
								Id:          3,
								Artist:      "Popeye",
								ResourceUrl: "http://somehting.com/release/3",
							},
						},
					},
				}
				callCount := 0
				handlerFunc := func(w http.ResponseWriter, _ *http.Request) {
					if callCount >= len(responses) {
						t.Fatalf("unexpected call to handler: %d", callCount)
					}

					resp := responses[callCount]
					callCount++
					w.WriteHeader(200)
					json.NewEncoder(w).Encode(resp)
				}
				return handlerFunc
			},
			expected: []Release{
				{
					Id:          1,
					Artist:      "John Doe",
					ResourceUrl: "http://somehting.com/release/1",
				},
				{
					Id:          2,
					Artist:      "Marry Poppins",
					ResourceUrl: "http://somehting.com/release/2",
				},
				{
					Id:          3,
					Artist:      "Popeye",
					ResourceUrl: "http://somehting.com/release/3",
				},
			},
		},
		{
			name: "if API response code is other than 200, the error is returned",
			handler: func() http.HandlerFunc {
				return func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(500)
					json.NewEncoder(w).Encode(ReleasesResponse{})
				}
			},
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler())
			httpClient := &http.Client{}
			client := &Client{
				HttpClient: httpClient,
			}

			actual, err := Releases(client, server.URL)

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
