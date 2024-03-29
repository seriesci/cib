package golang

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cli"
)

// Run runs all Go related stuff.
// At the moment only total test coverage.
func Run() error {
	// create code coverage file
	profileArgs := []string{
		"test",
		"./...",
		"-coverprofile",
		"cover.out",
	}
	profileCmd := exec.Command("go", profileArgs...)
	profileCmd.Stdout = os.Stdout
	profileCmd.Stderr = os.Stderr
	if err := profileCmd.Run(); err != nil {
		return err
	}

	cli.Checkf("cover profile %s created\n", cli.Blue("cover.out"))

	// read code coverage file
	funcArgs := []string{
		"tool",
		"cover",
		"-func",
		"cover.out",
	}
	cmd := exec.Command("go", funcArgs...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	var total string
	for scanner.Scan() {
		text := scanner.Text()

		// check for last line, e.g. "total: (statements) 86.7%"
		if !strings.HasPrefix(text, "total:") {
			continue
		}

		// get the total percentage from the last line
		fields := strings.Fields(text)
		total = fields[len(fields)-1]
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	cli.Checkln("total coverage (statements) is", cli.Blue(total))

	// create series
	if err := api.CreateSeries(api.SeriesCoverage); err != nil {
		return err
	}

	if err := api.Post(total, api.SeriesCoverage); err != nil {
		return err
	}

	// build binary to get size
	buildArgs := []string{
		"build",
		"-o",
		"binary_by_seriesci",
	}
	buildCmd := exec.Command("go", buildArgs...)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}

	cli.Checkf("binary %s built\n", cli.Blue("binary_by_seriesci"))

	info, err := os.Stat("binary_by_seriesci")
	if err != nil {
		return err
	}

	size := fmt.Sprintf("%.2fMB", float64(info.Size())/1000/1000)
	cli.Checkln("binary file size is", cli.Blue(size))

	// create series
	if err := api.CreateSeries(api.SeriesFileSize); err != nil {
		return err
	}

	if err := api.Post(size, api.SeriesFileSize); err != nil {
		return err
	}

	return nil
}
