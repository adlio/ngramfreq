package main

import (
	"strings"
	"testing"
)

func Test3WordString(t *testing.T) {
	r := strings.NewReader("my first trigram")
	results := extractNgrams(r)

	freq, _ := results["my first trigram"]
	if freq != 1 {
		t.Errorf("Expected 'my first trigram' freq of 1. Got %d", freq)
	}
}

func TestPunctuationCleanup(t *testing.T) {
	r := strings.NewReader(`
			I love
			sandwiches. Very much.
			Lorem ipsum dolor sit amet.
			I love sandwiches. (I LOVE SANDWICHES!!).
			East Side Deli makes the best sandwiches with love.
			`)

	results := extractNgrams(r)

	freq, _ := results["i love sandwiches"]
	if freq != 3 {
		t.Errorf("Expected 'i love sandwiches' freq of 3. Got %d", freq)
	}
}
