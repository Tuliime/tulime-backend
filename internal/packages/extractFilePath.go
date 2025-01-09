package packages

import (
	"fmt"
	"net/url"
	"strings"
)

// ExtractFilePath extracts the file path from a Firebase Storage URL
// and decodes URL-safe characters.
func ExtractFilePath(firebaseURL string) (string, error) {
	parsedURL, err := url.Parse(firebaseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	filePath := parsedURL.Path
	if !strings.HasPrefix(filePath, "/v0/b/") {
		return "", fmt.Errorf("unexpected URL format")
	}

	pathParts := strings.Split(filePath, "/o/")
	if len(pathParts) < 2 {
		return "", fmt.Errorf("file path not found in URL")
	}

	decodedPath, err := url.QueryUnescape(pathParts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode file path: %w", err)
	}

	return decodedPath, nil
}
