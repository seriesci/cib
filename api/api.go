package api

import (
	"bytes"
	"encoding/json"
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

// Descriptions returns the description for a given series name.
var Descriptions = map[string]string{
	SeriesCoverage:      "Code coverage",
	SeriesFileSize:      "File size",
	SeriesTime:          "Build time",
	SeriesBundleSize:    "Bundle size",
	SeriesDependencies:  "Number of dependencies",
	SeriesPerformance:   "Lighthouse performance",
	SeriesAccessibility: "Lighthouse accessibility",
	SeriesPractices:     "Lighthouse best practices",
	SeriesSEO:           "Lighthouse SEO",
}

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

// Repo returns owner and repository in the form owner/repository.
func Repo() (string, error) {
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
	r, err := Repo()
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

	cli.Checkf("post %s: status code: %s, body: %s\n", cli.Blue(series), cli.Blue(res.StatusCode), cli.Blue(string(body)))

	return nil
}

// CreateSeriesRequest request
type CreateSeriesRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateSeries creates a new series.
func CreateSeries(series string) error {

	// create custom request
	data := CreateSeriesRequest{
		Name:        series,
		Description: Descriptions[series],
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(data); err != nil {
		return err
	}

	// get repo in form owner/repo
	r, err := Repo()
	if err != nil {
		return err
	}

	u := fmt.Sprintf("https://seriesci.com/api/repos/%s/series", r)
	req, err := http.NewRequest(http.MethodPost, u, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("SERIESCI_TOKEN")))

	// send request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusConflict {
		cli.Checkf("series %s already exists\n", cli.Blue(series))
	} else {
		cli.Checkf("series %s created\n", cli.Blue(series))
	}

	return nil
}
