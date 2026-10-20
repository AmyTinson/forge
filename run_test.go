package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteNewFileDoesNotOverwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.txt")
	if err := os.WriteFile(path, []byte("original"), 0644); err != nil {
		t.Fatal(err)
	}

	err := writeNewFile(path, []byte("replacement"))
	if !errors.Is(err, os.ErrExist) {
		t.Fatalf("writeNewFile error = %v, want os.ErrExist", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "original" {
		t.Fatalf("file content = %q, want original content", data)
	}
}
