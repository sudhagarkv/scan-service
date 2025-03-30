package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func SplitGitHubURL(githubURL string) (string, string, error) {
	// Parse the GitHub URL
	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		return "", "", err
	}

	// Extract path from the parsed URL
	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")

	// If the path has at least two parts, consider the first part as the namespace and the second as the repo name
	if len(pathParts) >= 2 {
		namespace := pathParts[0]
		repoName := pathParts[1]

		return namespace, strings.TrimSuffix(repoName, ".git"), nil
	}

	// If the path does not have enough parts, return an error
	return "", "", fmt.Errorf("invalid GitHub URL format")
}
