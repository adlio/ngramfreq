package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

// GramSize is the target length of word
// sequences (i.e. 3 means we search for
// trigrams)
var GramSize = 3

// OutputSize is the number of n-grams to
// output
var OutputSize = 100

// Grams maps all encountered n-grams to their
// frequencies
var Grams = make(map[string]*NGramFreq)

// Freqs stores every encountered n-gram so that
// it can be sorted for scoring
var Freqs = make([]*NGramFreq, 0)

// main is the entry point for the application
func main() {

	flag.IntVar(&GramSize, "s", 3, "sequence size (2 = bigrams, 3=trigrams)")
	flag.IntVar(&OutputSize, "n", 100, "number of n-grams to output")
	flag.Parse()

	// Get filenames from arguments
	var filenames []string
	for _, arg := range flag.Args() {
		filename := path.Clean(arg)
		if info, err := os.Stat(filename); err == nil && info.Mode().IsRegular() {
			filenames = append(filenames, filename)
		} else {
			log.Printf("Ignoring invalid file '%s'.")
		}
	}

	// Check for input on STDIN
	var haveStdIn bool
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err) // TODO Improve error message
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		haveStdIn = true
	}

	if len(filenames) == 0 && !haveStdIn {
		invalidArgs()
	}

	for _, filename := range filenames {
		processFile(filename)
	}

	if haveStdIn {
		extractNgrams(os.Stdin)
	}

	sort.Slice(Freqs, func(i, j int) bool {
		return Freqs[i].Freq > Freqs[j].Freq
	})

	for i := 0; i < OutputSize && i < len(Freqs)-1; i++ {
		fmt.Println(Freqs[i])
	}

}

// processFile extracts all n-grams from the
// supplied file
//
// TODO Add readability checking
func processFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	extractNgrams(f)
}

// extractNgrams tokenizes the supplied stream and
// extracts each unique n-gram into a map relating
// the n-gram to its frequency
func extractNgrams(r io.Reader) {

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
		if len(wordQueue) > GramSize {
			wordQueue = wordQueue[1:] // Drop the first element
		}
		if len(wordQueue) == GramSize {
			phrase := strings.Join(wordQueue, " ")

			if ngf, existed := Grams[phrase]; existed {
				ngf.Freq++
			} else {
				ngf = &NGramFreq{Text: phrase, Freq: 1}
				Grams[phrase] = ngf
				Freqs = append(Freqs, ngf)
			}
		}
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

func invalidArgs() {
	fmt.Printf("Usage: %s [option flags] file1.txt file2.txt\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}
