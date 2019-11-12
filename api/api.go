package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/seriesci/cib/cli"
)

// series names
const (
	SeriesCoverage      = "coverage"
	SeriesFileSize      = "size"
	SeriesTime          = "time"
	SeriesBundleSize    = "bundlesize"
	SeriesDependencies  = "dependencies"
	SeriesPerformance   = "performance"
	SeriesAccessibility = "accessibility"
	SeriesPractices     = "practices"
	SeriesSEO           = "seo"
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

func repo() (string, error) {
	// github actions
	repo, ok := os.LookupEnv("GITHUB_REPOSITORY")
	if ok {
		return repo, nil
	}

	// travis ci
	repo, ok = os.LookupEnv("TRAVIS_REPO_SLUG")
	if ok {
		return repo, nil
	}

	// circle ci
	if _, ok = os.LookupEnv("CIRCLECI"); ok {
		username := os.Getenv("CIRCLE_PROJECT_USERNAME")
		reponame := os.Getenv("CIRCLE_PROJECT_REPONAME")
		return username + "/" + reponame, nil
	}

	return "", errors.New("cannot find repo in environment variables")
}

// Post something.
// currently only works with GitHub Actions.
func Post(value, series string) error {
	// get commit hash
	s, err := sha()
	if err != nil {
		return err
	}

	data := url.Values{
		"value": {value},
		"sha":   {s},
	}

	// get repo in form owner/repo
	r, err := repo()
	if err != nil {
		return err
	}

	u := fmt.Sprintf("https://seriesci.com/api/repos/%s/%s/combined", r, series)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("SERIESCI_TOKEN")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cli.Checkf("post %s: status code: %s, body: %s\n", series, res.StatusCode, string(body))

	return nil
}
