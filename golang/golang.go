package golang

import (
	"bufio"
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

	cli.Checkln("total coverage (statements) is", cli.Blue(total))

	res, err := api.Post(total)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cli.Checkf("response status code: %s, body: %s", res.StatusCode, string(body))

	// done
	cli.Checkln("I'm done. See", cli.Blue("https://seriesci.com/seriesci/cib/series/master/coverage"))

	return nil
}
