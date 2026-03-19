package httputil

import (
	"net/http"
	"time"
)

// Do a simple get http request and return a K response data.
func Get[K any](url string) (*K, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	var data K
	_, err = DoReq(client, req, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
