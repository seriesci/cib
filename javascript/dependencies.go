package javascript

import (
	"errors"
	"fmt"

	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cli"
)

func dependencies(packageJSON map[string]interface{}) error {
	dependencies, ok := packageJSON["dependencies"].(map[string]interface{})
	if !ok {
		return errors.New("dependencies not found")
	}

	cli.Checkf("%s dependencies found\n", blue(len(dependencies)))

	// create series
	if err := api.CreateSeries(api.SeriesDependencies); err != nil {
		return err
	}

	if err := api.Post(fmt.Sprintf("%d", len(dependencies)), api.SeriesDependencies); err != nil {
		return err
	}

	return nil
}
