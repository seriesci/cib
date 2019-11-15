package javascript

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cli"
	"github.com/seriesci/cib/lighthouse"
)

func runLighthouse() error {
	install := exec.Command("npm", "install", "--production", "lighthouse")
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	if err := install.Run(); err != nil {
		return err
	}

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
	args := []string{
		"lighthouse",
		"http://localhost:3000",
		"--output=json",
		"--output-path=./lighthouse.json",
		`--chrome-flags="--headless"`,
	}
	cmd := exec.Command("npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
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

	performance := fmt.Sprintf("%.2f%%", report.Categories.Performance.Score*100)
	cli.Checkf("lighthouse performance is %s\n", cli.Blue(performance))
	if err := api.CreateSeries(api.SeriesPerformance); err != nil {
		return err
	}
	if err := api.Post(performance, api.SeriesPerformance); err != nil {
		return err
	}

	accessibility := fmt.Sprintf("%.2f%%", report.Categories.Accessibility.Score*100)
	cli.Checkf("lighthouse accessibility is %s\n", cli.Blue(accessibility))
	if err := api.CreateSeries(api.SeriesAccessibility); err != nil {
		return err
	}
	if err := api.Post(accessibility, api.SeriesAccessibility); err != nil {
		return err
	}

	bestPractices := fmt.Sprintf("%.2f%%", report.Categories.BestPractices.Score*100)
	cli.Checkf("lighthouse best practices is %s\n", cli.Blue(bestPractices))
	if err := api.CreateSeries(api.SeriesPractices); err != nil {
		return err
	}
	if err := api.Post(bestPractices, api.SeriesPractices); err != nil {
		return err
	}

	seo := fmt.Sprintf("%.2f%%", report.Categories.Seo.Score*100)
	cli.Checkf("lighthouse seo is %s\n", cli.Blue(seo))
	if err := api.CreateSeries(api.SeriesSEO); err != nil {
		return err
	}
	if err := api.Post(seo, api.SeriesSEO); err != nil {
		return err
	}

	return nil
}
