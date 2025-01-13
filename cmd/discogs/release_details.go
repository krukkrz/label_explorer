package discogs

import (
	"encoding/json"
	"fmt"
)

type DetailsResponse struct {
	Id     int      `json:"id,omitempty"`
	Genres []string `json:"genres,omitempty"`
	Styles []string `json:"styles,omitempty"`
}

func ReleaseDetails(
	client *Client,
	url string,
) (DetailsResponse, error) {
	req, err := buildGetRequest(client, url)
	if err != nil {
		return DetailsResponse{}, fmt.Errorf("failed to build authorized request: %w", err)
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return DetailsResponse{}, fmt.Errorf("failed to call url %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return DetailsResponse{}, fmt.Errorf("incorrect response %d returned from url %s", resp.StatusCode, url)
	}

	var apiResp DetailsResponse
	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return DetailsResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}
	return apiResp, nil
}
