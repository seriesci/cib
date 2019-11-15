package javascript

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fatih/color"
)

const (
	check = "\u2713"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	blue  = color.New(color.FgBlue).SprintFunc()
)

// Run runs all JavaScript related stuff.
func Run() error {
	fmt.Printf("cib: %s language %s detected\n", green(check), blue("JavaScript"))

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

	fmt.Printf("cib: %s build script %s found\n", green(check), blue(build))

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
