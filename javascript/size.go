package javascript

import (
	"fmt"

	"github.com/seriesci/cib/api"
	"github.com/seriesci/cib/cli"
	"github.com/seriesci/cib/size"
)

func bundlesize() error {
	s, err := size.Directory("build")
	if err != nil {
		return err
	}

	str := fmt.Sprintf("%fK", s)
	cli.Checkf("total size of \"build\" directory is %s\n", blue(str))

	// create series
	if err := api.CreateSeries(api.SeriesBundleSize); err != nil {
		return err
	}

	// post value
	if err := api.Post(str, api.SeriesBundleSize); err != nil {
		return err
	}

	return nil
}
