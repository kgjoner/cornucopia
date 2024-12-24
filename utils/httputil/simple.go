package httputil

import (
	"net/http"
	"time"
)

// Do a simple get http request and return a K response data.
func Get[K any](url string) (data *K, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	_, err = DoReq(client, req, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
