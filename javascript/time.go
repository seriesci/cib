package javascript

import (
	"os"
	"os/exec"
	"time"

	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cli"
)

func duration() error {
	start := time.Now()

	cmd := exec.Command("npm", "run", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	elapsed := time.Since(start)

	cli.Checkf("build took %s\n", cli.Blue(elapsed))

	// create series
	if err := api.CreateSeries(api.SeriesTime); err != nil {
		return err
	}

	// post value
	if err := api.Post(elapsed.String(), api.SeriesTime); err != nil {
		return err
	}

	return nil
}
