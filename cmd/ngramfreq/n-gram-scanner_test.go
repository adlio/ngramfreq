package main

import (
	"strings"
	"testing"
)

func TestEmptyString(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader("")
	s.Scan(r)

	if len(s.Grams) != 0 {
		t.Error("An empty string shouldn't have any ngrams in it.")
	}
}

func Test3WordString(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader("my first trigram")
	s.Scan(r)

	ngf, _ := s.Grams["my first trigram"]
	if ngf.Freq != 1 {
		t.Errorf("Expected 'my first trigram' freq of 1. Got %d", ngf.Freq)
	}

	_, found := s.Grams["  my"]
	if found {
		t.Errorf("Grams contained something that wasn't a trigram.")
	}
}

func TestPunctuationOnlyString(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader("lorem ipsum dolor ... --- ... --- ... --- ... ---")
	s.Scan(r)

	ngf, _ := s.Grams["lorem ipsum dolor"]
	if ngf.Freq != 1 {
		t.Errorf("Incorrect .Freq on 'lorem ipsum dolor': %d", ngf.Freq)
	}

	ngf, found := s.Grams["  "]
	if found {
		t.Errorf("All-spaces n-gram found '%s'", ngf)
	}
}

func TestPunctuationCleanup(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader(`
			I love
			sandwiches. Very much.
			Lorem ipsum dolor sit amet. Chapter 10
			I love sandwiches. (I LOVE SANDWICHES!!).
			East Side Deli makes the best sandwiches with love.
			`)
	s.Scan(r)

	ngf, ok := s.Grams["i love sandwiches"]
	if ok != true {
		t.Fatal("Should have found an ngram for 'i love sandwiches'")
	}
	if ngf.Freq != 3 {
		t.Errorf("Expected 'i love sandwiches' freq of 3. Got %d", ngf.Freq)
	}
}

func TestScrubWhenClean(t *testing.T) {
	s := ScrubWord("cleanword")
	if s != "cleanword" {
		t.Errorf("Didn't need scrubbing, but changed to '%s'", s)
	}
}

func TestScrubWordLowerCasing(t *testing.T) {
	s := ScrubWord("CleanWord")
	if s != "cleanword" {
		t.Errorf("Expected 'cleanword', got '%s'.", s)
	}
}
