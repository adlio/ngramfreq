package main

import "fmt"

// NGramFreq stores an ngram and its frequency.
//
type NGramFreq struct {
	Text string
	Freq int64
}

// String formats an NGramFreq for console output
func (n *NGramFreq) String() string {
	return fmt.Sprintf("%d - %s", n.Freq, n.Text)
}
