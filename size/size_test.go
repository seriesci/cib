package size

import "testing"

func TestDirectory(t *testing.T) {
	s, err := Directory("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if s != 123 {
		t.Fatalf("got %d; want %d", s, 123)
	}
}
