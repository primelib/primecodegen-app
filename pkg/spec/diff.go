package spec

import (
	"fmt"
	"os"
)

type Diff struct {
	OpenAPI []OpenAPIDiff
}

func DiffSpec(format string, file1 string, file2 string) (Diff, error) {
	var diff Diff

	// check of files exist
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		return diff, fmt.Errorf("file %s does not exist", file1)
	}
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		return diff, fmt.Errorf("file %s does not exist", file2)
	}

	// diff openapi
	if format == "openapi" {
		d, err := DiffOpenAPI(file1, file2)
		if err != nil {
			return diff, fmt.Errorf("failed to diff openapi: %w", err)
		}
		diff.OpenAPI = d
	}

	return diff, nil
}
