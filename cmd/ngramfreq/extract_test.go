package main

import (
	"strings"
	"testing"
)

func TestNGramFreqString(t *testing.T) {
	ngf := &NGramFreq{Text: "message number one", Freq: 45}
	s := ngf.String()
	expected := "45 - message number one"
	if s != expected {
		t.Errorf("Expected '%s', got '%s'", expected, s)
	}
}

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

	_, found := Grams["  my"]
	if found {
		t.Errorf("Grams contained something that wasn't a trigram.")
	}
}

func TestPunctuationOnlyString(t *testing.T) {
	Grams = make(map[string]*NGramFreq)
	r := strings.NewReader("lorem ipsum dolor ... --- ... --- ... --- ... ---")
	extractNgrams(r)

	ngf, _ := Grams["lorem ipsum dolor"]
	if ngf.Freq != 1 {
		t.Errorf("Incorrect .Freq on 'lorem ipsum dolor': %d", ngf.Freq)
	}

	ngf, found := Grams["  "]
	if found {
		t.Errorf("All-spaces n-gram found '%s'", ngf)
	}
}

func TestPunctuationCleanup(t *testing.T) {
	Grams = make(map[string]*NGramFreq)
	r := strings.NewReader(`
			I love
			sandwiches. Very much.
			Lorem ipsum dolor sit amet. Chapter 10
			I love sandwiches. (I LOVE SANDWICHES!!).
			East Side Deli makes the best sandwiches with love.
			`)

	extractNgrams(r)

	ngf, ok := Grams["i love sandwiches"]
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
