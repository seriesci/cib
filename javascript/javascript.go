package javascript

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/seriesci/cib/cli"
)

// Run runs all JavaScript related stuff.
func Run() error {
	cli.Checkf("language %s detected\n", cli.Blue("JavaScript"))

	// check if node_modules folder exists
	if _, err := os.Stat("node_modules"); os.IsNotExist(err) {
		cli.Checkf("could not find %s. running %s\n", cli.Blue("node_modules"), cli.Blue("npm ci"))

		if err := install(); err != nil {
			return err
		}
	}

	// run build script
	packageJSON, err := ioutil.ReadFile("package.json")
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(packageJSON), &result); err != nil {
		return err
	}

	scripts, ok := result["scripts"].(map[string]interface{})
	if !ok {
		return errors.New("scripts not found")
	}

	build, ok := scripts["build"]
	if !ok {
		return errors.New("build script not found")
	}

	cli.Checkf("build script %s found\n", cli.Blue(build))

	// run the build
	if err := duration(); err != nil {
		return err
	}

	// run bundle size
	if err := bundlesize(); err != nil {
		return err
	}

	// get code coverage
	if err := coverage(result); err != nil {
		return err
	}

	// count dependencies
	if err := dependencies(result); err != nil {
		return err
	}

	// run lighthouse
	if err := runLighthouse(); err != nil {
		return err
	}

	return nil
}
