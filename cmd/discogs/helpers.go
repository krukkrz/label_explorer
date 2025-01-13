package discogs

import (
	"fmt"
	"net/http"
)

func buildGetRequest(client *Client, url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate a request: %w", err)
	}
	req.Header["Authorization"] = []string{fmt.Sprintf("Discogs key=%s, secret=%s", client.Key, client.Secret)}
	return req, nil
}
