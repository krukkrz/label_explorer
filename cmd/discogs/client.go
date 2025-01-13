package discogs

import "net/http"

type Client struct {
	HttpClient *http.Client
	Key        string
	Secret     string
}
