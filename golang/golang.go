package golang

import (
	"bufio"
	"fmt"
	"io/ioutil"
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

	cli.Checkf("cover profile %s created\n", "cover.out")

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

	cli.Checkln("total coverage (statements) is", total)

	res, err := api.Post(total, api.SeriesCoverage)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cli.Checkf("post %s: res status code: %s, body: %s\n", api.SeriesCoverage, res.StatusCode, string(body))

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

	cli.Checkf("binary %s built\n", "binary_by_seriesci")

	info, err := os.Stat("binary_by_seriesci")
	if err != nil {
		return err
	}

	size := fmt.Sprintf("%.2fMB", float64(info.Size())/1000/1000)
	cli.Checkln("binary file size is", size)

	res, err = api.Post(size, api.SeriesFileSize)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cli.Checkf("post %s: res status code: %s, body: %s\n", api.SeriesFileSize, res.StatusCode, string(body))

	// done
	cli.Checkln("I'm done. See", "https://seriesci.com/seriesci/cib")

	return nil
}
