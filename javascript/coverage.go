package javascript

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cli"
	"github.com/seriesci/cib/cov"
)

func coverage(packageJSON map[string]interface{}) error {
	// edit package.json and add clover coverage reporter
	packageJSON["jest"] = map[string][]string{
		"coverageReporters": []string{"clover"},
	}

	// override package.json temporarily
	b, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("package.json", b, 0644); err != nil {
		return err
	}

	// run code coverage
	covArgs := []string{
		"test",
		"--",
		"--coverage",
		"--watchAll=false",
	}
	covCmd := exec.Command("npm", covArgs...)
	covCmd.Stdout = os.Stdout
	covCmd.Stderr = os.Stderr
	if err := covCmd.Run(); err != nil {
		return err
	}

	// read coverage
	cloverXML, err := ioutil.ReadFile(filepath.Join("coverage", "clover.xml"))
	if err != nil {
		return err
	}

	var coverage cov.Coverage
	if err := xml.Unmarshal(cloverXML, &coverage); err != nil {
		return err
	}

	// covered statements
	cs := float64(coverage.Project.Metrics.CoveredStatements) / float64(coverage.Project.Metrics.Statements) * 100
	str := fmt.Sprintf("%.2f%%", cs)

	cli.Checkf("code coverage is %s\n", blue(str))

	// create series
	if err := api.CreateSeries(api.SeriesCoverage); err != nil {
		return err
	}

	// post value
	if err := api.Post(str, api.SeriesCoverage); err != nil {
		return err
	}

	return nil

}
