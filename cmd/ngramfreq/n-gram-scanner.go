package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Scanner struct {
	GramSize int
	Grams    map[string]*NGramFreq
	Freqs    []*NGramFreq
}

func NewScanner(gramSize int) *Scanner {
	return &Scanner{
		GramSize: gramSize,
		Grams:    make(map[string]*NGramFreq),
		Freqs:    make([]*NGramFreq, 0),
	}
}

func (s *Scanner) ScanFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	s.Scan(f)
	return nil
}

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

		wordQueue = append(wordQueue, word)
		if len(wordQueue) > s.GramSize {
			wordQueue = wordQueue[1:] // Drop the first element
		}
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

func (s *Scanner) Score() {
	sort.Slice(s.Freqs, func(i, j int) bool {
		return s.Freqs[i].Freq > s.Freqs[j].Freq
	})
}

func (s *Scanner) WriteTopN(n int, w io.Writer) {
	for i := 0; i < n && i < len(s.Freqs); i++ {
		fmt.Fprintln(w, s.Freqs[i])
	}
}
