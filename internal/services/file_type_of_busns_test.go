package services

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileTypeOfBusnsPreservesHyphenInDescription(t *testing.T) {
	file := filepath.Join(t.TempDir(), "vid.txt")
	if err := os.WriteFile(file, []byte("12345 - Description - with hyphen"), 0600); err != nil {
		t.Fatal(err)
	}

	types, err := NewFileTypeOfBusns(file)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := types["12345"], "Description - with hyphen"; got != want {
		t.Fatalf("description = %q, want %q", got, want)
	}
}
