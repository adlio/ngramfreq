package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Scanner provides functionality to scan and count
// unique n-grams in text. Use NewScanner to get a
// properly initialized instance.
//
// NOTE: Scanner internally relies on a non thread-safe
// map. Do not call functions on Scanner concurrently.
type Scanner struct {
	GramSize int
	Grams    map[string]*NGramFreq
	Freqs    []*NGramFreq
}

// NewScanner creates a Scanner targeting n-grams of the
// specified size.
func NewScanner(gramSize int) *Scanner {
	return &Scanner{
		GramSize: gramSize,
		Grams:    make(map[string]*NGramFreq),
		Freqs:    make([]*NGramFreq, 0),
	}
}

// ScanFile is a wrapper around Scan which accepts
// a file path instead of an already-opened file.
//
func (s *Scanner) ScanFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	s.Scan(f)
	return nil
}

// Scan iterates over the supplied io.Reader and collects
// n-gram results in Scanner.Grams and Scanner.Freqs.
//
func (s *Scanner) Scan(r io.Reader) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	wordQueue := make([]string, 0)

	for scanner.Scan() {
		word := ScrubWord(scanner.Text())

		// Skip words which scrubbing caused to be empty
		if word == "" {
			continue
		}

		// Pass each word through a fixed-size queue
		wordQueue = append(wordQueue, word)
		if len(wordQueue) > s.GramSize {
			wordQueue = wordQueue[1:] // Drop the first element
		}

		// Each time we get a full queue, process the n-gram.
		// Data is stored in pointers to NGramFreq so that
		// modifications via the map are reflected in the slice
		// as well.
		if len(wordQueue) == s.GramSize {
			phrase := strings.Join(wordQueue, " ")

			if ngf, existed := s.Grams[phrase]; existed {
				ngf.Freq++
			} else {
				ngf = &NGramFreq{Text: phrase, Freq: 1}
				s.Grams[phrase] = ngf
				s.Freqs = append(s.Freqs, ngf)
			}
		}
	}
}

// Score triggers a sort of .Freqs so that most-common
// n-grams are at the front of the slice.
func (s *Scanner) Score() {
	sort.Slice(s.Freqs, func(i, j int) bool {
		return s.Freqs[i].Freq > s.Freqs[j].Freq
	})
}

// WriteTopN outputs a "leaderboard" to the supplied
// writer.
func (s *Scanner) WriteTopN(n int, w io.Writer) {
	for i := 0; i < n && i < len(s.Freqs); i++ {
		fmt.Fprintln(w, s.Freqs[i])
	}
}

// ScrubWord scrubs noise characters (punctuation, etc)
// and lowercases the input
func ScrubWord(s string) string {

	// return strings.ToLower(invalidChars.ReplaceAllString(s, ""))
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r + 32
		case r >= '0' && r <= '9':
			return r
		default:
			return -1
		}
	}, s)
}
