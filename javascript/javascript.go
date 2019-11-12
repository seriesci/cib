package javascript

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
	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cov"
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
	start := time.Now()

	buildCmd := exec.Command("npm", "run", "build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}

	elapsed := time.Since(start)

	fmt.Printf("cib: %s build took %s\n", green(check), blue(elapsed))

	if err := api.Post(elapsed.String(), api.SeriesTime); err != nil {
		return err
	}

	// run bundle size
	s, err := size.Directory("build")
	if err != nil {
		return err
	}
	fmt.Printf("cib: %s total size of \"build\" directory is %s\n", green(check), blue(s, "kB"))

	if err := api.Post(fmt.Sprintf("%fK", s), api.SeriesBundleSize); err != nil {
		return err
	}

	// edit package.json and add clover coverage reporter
	result["jest"] = map[string][]string{
		"coverageReporters": []string{"clover"},
	}

	// override package.json temporarily
	b, err := json.MarshalIndent(result, "", "  ")
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

	fmt.Printf("cib: %s code coverage is %s\n", green(check), blue(fmt.Sprintf("%.2f%%", cs)))

	if err := api.Post(fmt.Sprintf("%.2f%%", cs), api.SeriesCoverage); err != nil {
		return err
	}

	// count dependencies
	dependencies, ok := result["dependencies"].(map[string]interface{})
	if !ok {
		return errors.New("dependencies not found")
	}

	fmt.Printf("cib: %s %s dependencies found\n", green(check), blue(len(dependencies)))

	if err := api.Post(fmt.Sprintf("%d", len(dependencies)), api.SeriesDependencies); err != nil {
		return err
	}

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
	lighthouseCMD.Stdout = os.Stdout
	lighthouseCMD.Stderr = os.Stderr
	if err := lighthouseCMD.Run(); err != nil {
		return err
	}

	if err := server.Shutdown(context.Background()); err != nil {
		return err
	}

	lighthouseJSON, err := ioutil.ReadFile("lighthouse.json")
	if err != nil {
		return err
	}

	var report lighthouse.Report
	if err := json.Unmarshal(lighthouseJSON, &report); err != nil {
		return err
	}

	performance := report.Categories.Performance.Score * 100
	accessibility := report.Categories.Accessibility.Score * 100
	bestPractices := report.Categories.BestPractices.Score * 100
	seo := report.Categories.Seo.Score * 100

	fmt.Printf("cib: %s lighthouse performance is %s\n", green(check), blue(performance, "%"))
	if err := api.Post(fmt.Sprintf("%.2f%%", performance), api.SeriesPerformance); err != nil {
		return err
	}

	fmt.Printf("cib: %s lighthouse accessibility is %s\n", green(check), blue(accessibility, "%"))
	if err := api.Post(fmt.Sprintf("%.2f%%", accessibility), api.SeriesAccessibility); err != nil {
		return err
	}

	fmt.Printf("cib: %s lighthouse best practices is %s\n", green(check), blue(bestPractices, "%"))
	if err := api.Post(fmt.Sprintf("%.2f%%", bestPractices), api.SeriesPractices); err != nil {
		return err
	}

	fmt.Printf("cib: %s lighthouse seo is %s\n", green(check), blue(seo, "%"))
	if err := api.Post(fmt.Sprintf("%.2f%%", seo), api.SeriesSEO); err != nil {
		return err
	}

	return nil
}
