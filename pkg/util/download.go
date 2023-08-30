package util

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, outputFilePath string) error {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Copy the response body to the output file
	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadString(url string) (string, error) {
	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// return response body as string
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
