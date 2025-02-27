package util

import (
	"io"
	"net/http"
)

func DownloadString(url string) ([]byte, error) {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// return response body as string
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
