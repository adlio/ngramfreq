package main

import (
	"os"
	"testing"
)

func TestArgParsingWithFakeFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"ngramfreq", "fake.file"}

	_, _, err := processArgs()
	if err == nil {
		t.Errorf("Expected fake.file to throw an error. It didn't.")
	}
}

func TestArgParsingWithValidFile(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"ngramfreq", "./testdata/sonnets.txt"}

	hasStdin, filenames, err := processArgs()
	if err != nil {
		t.Error(err)
	}
	if hasStdin != false {
		t.Errorf("hasStdin should be false. Got %v", hasStdin)
	}
	if len(filenames) != 1 {
		t.Errorf("Expected a single file with sonnets.txt. Got %v", filenames)
	}
}
