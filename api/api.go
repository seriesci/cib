package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// series names
const (
	SeriesCoverage = "coverage"
	SeriesFileSize = "size"
)

// Post something.
// currently only works with GitHub Actions.
func Post(value, series string) (*http.Response, error) {
	data := url.Values{
		"value": {value},
		"sha":   {os.Getenv("GITHUB_SHA")},
	}

	// IMPORTANT: this only works for GitHub Actions at the moment
	// seriesci/cib
	repo := os.Getenv("GITHUB_REPOSITORY")

	u := fmt.Sprintf("https://seriesci.com/api/repos/%s/%s/combined", repo, series)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	token := os.Getenv("SERIESCI_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return http.DefaultClient.Do(req)
}
