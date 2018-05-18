package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSonnets(t *testing.T) {
	s := NewScanner(3)
	err := s.ScanFile("./testdata/sonnets.txt")
	if err != nil {
		t.Error(err)
	}
	s.Score()

	assertScannerOutputPrefix(t, s, 100, "3 - thy self thy\n3 - why dost thou")
}

func TestScanNonexistentFile(t *testing.T) {
	s := NewScanner(3)
	err := s.ScanFile("./testdata/fakefile.txt")
	if err == nil {
		t.Error("./testdata/fakefile.txt doesn't exist. ScanFile() should have thrown an error about that.")
	}
}

func TestEmptyString(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader("")
	s.Scan(r)
	s.Score()

	if len(s.Grams) != 0 {
		t.Error("An empty string shouldn't have any ngrams.")
	}

	if len(s.Freqs) != 0 {
		t.Error("An empty string shouldn't have any ngrams.")
	}

	b := bytes.Buffer{}
	s.WriteTopN(100, &b)
	if b.String() != "" {
		t.Errorf("Empty input should result in empty output. Got: %v", b)
	}

}

func Test3WordString(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader("my first trigram")
	s.Scan(r)
	s.Score()

	ngf, _ := s.Grams["my first trigram"]
	if ngf.Freq != 1 {
		t.Errorf("Expected 'my first trigram' freq of 1. Got %d", ngf.Freq)
	}

	_, found := s.Grams["  my"]
	if found {
		t.Errorf("Grams contained something that wasn't a trigram.")
	}

	assertScannerOutputPrefix(t, s, 100, "1 - my first trigram")
}

func TestPunctuationOnlyString(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader("lorem ipsum dolor ... --- ... --- ... --- ... ---")
	s.Scan(r)
	s.Score()

	ngf, _ := s.Grams["lorem ipsum dolor"]
	if ngf.Freq != 1 {
		t.Errorf("Incorrect .Freq on 'lorem ipsum dolor': %d", ngf.Freq)
	}

	ngf, found := s.Grams["  "]
	if found {
		t.Errorf("All-spaces n-gram found '%s'", ngf)
	}

	assertScannerOutputPrefix(t, s, 100, "1 - lorem ipsum dolor")
}

func TestPunctuationCleanup(t *testing.T) {
	s := NewScanner(3)
	r := strings.NewReader(`
			It's important that I explain that I
			love
			sandwiches. Very much.
			Lorem ipsum dolor sit amet. Chapter 10
			I love sandwiches. (I LOVE SANDWICHES!!).
			East Side Deli makes the best sandwiches with love.
			`)
	s.Scan(r)
	s.Score()

	ngf, ok := s.Grams["i love sandwiches"]
	if ok != true {
		t.Fatal("Should have found an ngram for 'i love sandwiches'")
	}
	if ngf.Freq != 3 {
		t.Errorf("Expected 'i love sandwiches' freq of 3. Got %d", ngf.Freq)
	}

	assertScannerOutputPrefix(t, s, 100, "3 - i love sandwiches")
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

func assertScannerOutputPrefix(t *testing.T, s *Scanner, n int, prefix string) {
	b := bytes.Buffer{}
	s.WriteTopN(n, &b)
	if !strings.HasPrefix(b.String(), prefix) {
		t.Errorf("Expected output to start with\n\n%s\n\nActual output:\n%s", prefix, b.String())
	}
}
