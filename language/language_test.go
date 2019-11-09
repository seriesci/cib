package language

import (
	"path/filepath"
	"testing"
)

func TestUndefined(t *testing.T) {
	if _, err := Detect(filepath.Join("testdata", "error")); err == nil {
		t.Error("want err != nil")
	}
}

func TestJavaScript(t *testing.T) {
	lang, err := Detect(filepath.Join("testdata", "javascript"))
	if err != nil {
		t.Fatal(err)
	}
	if lang != JavaScript {
		t.Fatalf("got %d; want %d", lang, JavaScript)
	}
}

func TestGo(t *testing.T) {
	lang, err := Detect(filepath.Join("testdata", "go"))
	if err != nil {
		t.Fatal(err)
	}
	if lang != Go {
		t.Fatalf("got %d; want %d", lang, Go)
	}
}
