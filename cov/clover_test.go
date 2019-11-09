package cov

import (
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestStatements(t *testing.T) {
	b, err := ioutil.ReadFile(filepath.Join("testdata", "clover.xml"))
	if err != nil {
		t.Fatal(err)
	}

	var c Coverage
	if err := xml.Unmarshal(b, &c); err != nil {
		t.Fatal(err)
	}

	if c.Project.Metrics.Statements != 43 {
		t.Fatalf("got %d; want %d", c.Project.Metrics.Statements, 43)
	}

	if c.Project.Metrics.CoveredStatements != 2 {
		t.Fatalf("got %d; want %d", c.Project.Metrics.CoveredStatements, 2)
	}
}
