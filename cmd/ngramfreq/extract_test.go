package main

import (
	"strings"
	"testing"
)

func TestExtractNgrams(t *testing.T) {
	r := strings.NewReader("my first trigram")
	results := extractNgrams(r)

	freq, _ := results["my first trigram"]
	if freq != 1 {
		t.Errorf("Expected 'my first trigram' freq of 1. Got %d", freq)
	}
}
