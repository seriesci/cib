package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"

	"github.com/seriesci/cib/cov"
	"github.com/seriesci/cib/language"
	"github.com/seriesci/cib/lighthouse"
	"github.com/seriesci/cib/size"
)

const (
	check = "\u2713"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	blue  = color.New(color.FgBlue).SprintFunc()
)

func main() {
	// check token
	_, ok := os.LookupEnv("SERIESCI_TOKEN")
	if !ok {
		panic(errors.New("cannot find SERIESCI_TOKEN environment variable"))
	}
	fmt.Printf("cib: %s environment variable %s found\n", green(check), blue("SERIESCI_TOKEN"))

	// check programming language
	_, err := language.Detect(".")
	if err != nil {
		panic(err)
	}

	fmt.Printf("cib: %s language %s detected\n", green(check), blue("JavaScript"))

	// run build script
	packageJSON, err := ioutil.ReadFile("package.json")
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(packageJSON), &result); err != nil {
		panic(err)
	}

	scripts, ok := result["scripts"].(map[string]interface{})
	if !ok {
		panic(errors.New("scripts not found"))
	}

	build, ok := scripts["build"]
	if !ok {
		panic(errors.New("build script not found"))
	}

	fmt.Printf("cib: %s build script %s found\n", green(check), blue(build))

	// run the build
	start := time.Now()

	buildCmd := exec.Command("npm", "run", "build")
	// buildCmd.Stdout = os.Stdout
	// buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		panic(err)
	}

	elapsed := time.Since(start)

	fmt.Printf("cib: %s build took %s\n", green(check), blue(elapsed))

	// run bundle size
	s, err := size.Directory("build")
	if err != nil {
		panic(err)
	}
	fmt.Printf("cib: %s total size of \"build\" directory is %s\n", green(check), blue(s, "kB"))

	// edit package.json and add clover coverage reporter
	result["jest"] = map[string][]string{
		"coverageReporters": []string{"clover"},
	}

	// override package.json temporarily
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("package.json", b, 0644); err != nil {
		panic(err)
	}

	// run code coverage
	covArgs := []string{
		"test",
		"--",
		"--coverage",
		"--watchAll=false",
	}
	covCmd := exec.Command("npm", covArgs...)
	// covCmd.Stdout = os.Stdout
	// covCmd.Stderr = os.Stderr
	if err := covCmd.Run(); err != nil {
		panic(err)
	}

	// read coverage
	cloverXML, err := ioutil.ReadFile(filepath.Join("coverage", "clover.xml"))
	if err != nil {
		panic(err)
	}

	var coverage cov.Coverage
	if err := xml.Unmarshal(cloverXML, &coverage); err != nil {
		panic(err)
	}

	// covered statements
	cs := float64(coverage.Project.Metrics.CoveredStatements) / float64(coverage.Project.Metrics.Statements) * 100

	fmt.Printf("cib: %s code coverage is %s\n", green(check), blue(fmt.Sprintf("%.2f%%", cs)))

	// count dependencies
	dependencies, ok := result["dependencies"].(map[string]interface{})
	if !ok {
		panic(errors.New("dependencies not found"))
	}

	fmt.Printf("cib: %s %s dependencies found\n", green(check), blue(len(dependencies)))

	// run lighthouse
	http.Handle("/", http.FileServer(http.Dir("build")))

	server := &http.Server{
		Addr:    ":3000",
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// start lighthouse
	lighthouseArgs := []string{
		"lighthouse",
		"http://localhost:3000",
		"--output=json",
		"--output-path=./lighthouse.json",
		`--chrome-flags="--headless"`,
	}
	lighthouseCMD := exec.Command("npx", lighthouseArgs...)
	// lighthouseCMD.Stdout = os.Stdout
	// lighthouseCMD.Stderr = os.Stderr
	if err := lighthouseCMD.Run(); err != nil {
		panic(err)
	}

	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	lighthouseJSON, err := ioutil.ReadFile("lighthouse.json")
	if err != nil {
		panic(err)
	}

	var report lighthouse.Report
	if err := json.Unmarshal(lighthouseJSON, &report); err != nil {
		panic(err)
	}

	performance := report.Categories.Performance.Score * 100
	accessibility := report.Categories.Accessibility.Score * 100
	bestPractices := report.Categories.BestPractices.Score * 100
	seo := report.Categories.Seo.Score * 100

	fmt.Printf("cib: %s lighthouse performance is %s\n", green(check), blue(performance, "%"))
	fmt.Printf("cib: %s lighthouse accessibility is %s\n", green(check), blue(accessibility, "%"))
	fmt.Printf("cib: %s lighthouse best practices is %s\n", green(check), blue(bestPractices, "%"))
	fmt.Printf("cib: %s lighthouse seo is %s\n", green(check), blue(seo, "%"))
}
