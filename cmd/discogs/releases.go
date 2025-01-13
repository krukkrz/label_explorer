package discogs

import (
	"encoding/json"
	"fmt"
)

type ReleasesResponse struct {
	Pagination Pagination `json:"pagination"`
	Releases   []Release  `json:"releases,omitempty"`
}

type Pagination struct {
	Page    int `json:"page,omitempty"`
	Pages   int `json:"pages,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Items   int `json:"items,omitempty"`
}

type Release struct {
	Id          int    `json:"id,omitempty"`
	Artist      string `json:"artist,omitempty"`
	ResourceUrl string `json:"resource_url,omitempty"`
}

func Releases(
	client *Client,
	baseUrl string,
) ([]Release, error) {
	releases, err := fetchAllPages(client, baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases for label: %w", err)
	}
	return releases, nil
}

func fetchAllPages(
	client *Client,
	apiURL string,
) ([]Release, error) {
	var allItems []Release
	page := 1
	limit := 100 // discogs limit: https://www.discogs.com/developers#page:home,header:home-pagination

	for {
		apiResp, err := fetchPage(client, apiURL, page, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		allItems = append(allItems, apiResp.Releases...)

		if apiResp.Pagination.Page == apiResp.Pagination.Pages {
			break
		}

		// Update the page number for the next request
		page++
	}

	return allItems, nil
}

func fetchPage(
	client *Client,
	apiURL string,
	page int,
	limit int,
) (ReleasesResponse, error) {
	url := fmt.Sprintf("%s?page=%d&per_page=%d", apiURL, page, limit)
	req, err := buildGetRequest(client, url)
	if err != nil {
		return ReleasesResponse{}, fmt.Errorf("failed to build authorized request: %w", err)
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return ReleasesResponse{}, fmt.Errorf("failed to call url %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ReleasesResponse{}, fmt.Errorf("incorrect response %d returned from url %s", resp.StatusCode, url)
	}

	var apiResp ReleasesResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return ReleasesResponse{}, fmt.Errorf("failed to decode response for page %d: %w", page, err)
	}
	return apiResp, nil
}
