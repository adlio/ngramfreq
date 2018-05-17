package main

import (
	"strings"
	"testing"
)

func TestEmptyString(t *testing.T) {
	Grams = make(map[string]*NGramFreq)
	r := strings.NewReader("")
	extractNgrams(r)

	if len(Grams) != 0 {
		t.Error("An empty string shouldn't have any ngrams in it.")
	}
}

func Test3WordString(t *testing.T) {
	Grams = make(map[string]*NGramFreq)
	r := strings.NewReader("my first trigram")
	extractNgrams(r)

	ngf, _ := Grams["my first trigram"]
	if ngf.Freq != 1 {
		t.Errorf("Expected 'my first trigram' freq of 1. Got %d", ngf.Freq)
	}
}

func TestPunctuationCleanup(t *testing.T) {
	Grams = make(map[string]*NGramFreq)
	r := strings.NewReader(`
			I love
			sandwiches. Very much.
			Lorem ipsum dolor sit amet.
			I love sandwiches. (I LOVE SANDWICHES!!).
			East Side Deli makes the best sandwiches with love.
			`)

	extractNgrams(r)

	ngf, _ := Grams["i love sandwiches"]
	if ngf.Freq != 3 {
		t.Errorf("Expected 'i love sandwiches' freq of 3. Got %d", ngf.Freq)
	}
}
