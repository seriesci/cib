package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// series names
const (
	SeriesCoverage = "coverage"
	SeriesFileSize = "size"
)

// get sha depending on ci environment.
func sha() (string, error) {
	// github actions
	sha, ok := os.LookupEnv("GITHUB_SHA")
	if ok {
		return sha, nil
	}

	// travis ci
	sha, ok = os.LookupEnv("TRAVIS_COMMIT")
	if ok {
		return sha, nil
	}

	// circle ci
	sha, ok = os.LookupEnv("CIRCLE_SHA1")
	if ok {
		return sha, nil
	}

	// default git
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Post something.
// currently only works with GitHub Actions.
func Post(value, series string) (*http.Response, error) {
	s, err := sha()
	if err != nil {
		return nil, err
	}

	data := url.Values{
		"value": {value},
		"sha":   {s},
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
