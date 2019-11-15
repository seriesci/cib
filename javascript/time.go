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

	buildCmd := exec.Command("npm", "run", "build")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return err
	}

	elapsed := time.Since(start)

	cli.Checkf("build took %s\n", blue(elapsed))

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
