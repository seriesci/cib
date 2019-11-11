package lighthouse

import (
	"path/filepath"
	"testing"
)

func TestCategories(t *testing.T) {
	report, err := ReadUnmarshal(filepath.Join("testdata", "lighthouse.json"))
	if err != nil {
		t.Fatal(err)
	}

	if report.Categories.Performance.Score != 0.980000 {
		t.Fatalf("got %f; want %f", report.Categories.Performance.Score, 0.980000)
	}

	if report.Categories.Accessibility.Score != 0.940000 {
		t.Fatalf("got %f; want %f", report.Categories.Accessibility.Score, 0.940000)
	}

	if report.Categories.BestPractices.Score != 0.930000 {
		t.Fatalf("got %f; want %f", report.Categories.BestPractices.Score, 0.930000)
	}

	if report.Categories.Seo.Score != 0.920000 {
		t.Fatalf("got %f; want %f", report.Categories.Seo.Score, 0.920000)
	}
}
