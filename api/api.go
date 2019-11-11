package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Post something.
// currently only works with GitHub Actions.
func Post(value string) (*http.Response, error) {
	data := url.Values{
		"value": {value},
		"sha":   {os.Getenv("GITHUB_SHA")},
	}

	// IMPORTANT: this only works for GitHub Actions at the moment
	repo := os.Getenv("GITHUB_REPOSITORY")

	// todo: fix hardcoded series name "coverage"
	u := fmt.Sprintf("https://seriesci.com/api/repos/%s/coverage/combined", repo)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	token := os.Getenv("SERIESCI_TOKEN")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	return http.DefaultClient.Do(req)
}
