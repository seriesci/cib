package cov

import (
	"path/filepath"
	"testing"
)

func TestStatements(t *testing.T) {
	c, err := ReadUnmarshal(filepath.Join("testdata", "clover.xml"))
	if err != nil {
		t.Fatal(err)
	}

	if c.Project.Metrics.Statements != 43 {
		t.Fatalf("got %d; want %d", c.Project.Metrics.Statements, 43)
	}

	if c.Project.Metrics.CoveredStatements != 2 {
		t.Fatalf("got %d; want %d", c.Project.Metrics.CoveredStatements, 2)
	}
}
