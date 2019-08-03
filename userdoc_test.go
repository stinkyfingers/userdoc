package main

import (
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	filename := filepath.Join("test_files", "foobar.go")
	userdocs, err := parseFile(filename, []string{})
	if err != nil {
		t.Error(err)
	}
	t.Log("usersdocs: ", userdocs)
}
