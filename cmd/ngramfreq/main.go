package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

const GRAMSIZE = 3
const TOP_N = 100

var Grams = make(map[string]*NGramFreq)
var Freqs = make([]*NGramFreq, 0)

// NGramFreq stores an ngram and its frequency
//
type NGramFreq struct {
	Text string
	Freq int64
}

// String formats an NGramFreq for console output
func (n *NGramFreq) String() string {
	return fmt.Sprintf("%d - %s", n.Freq, n.Text)
}

func main() {
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

	fmt.Printf("Processing %d files: %v. StdIn: %v\n", len(filenames), filenames, haveStdIn)

	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		extractNgrams(f)
		f.Close()
	}

	if haveStdIn {
		extractNgrams(os.Stdin)
	}

	sort.Slice(Freqs, func(i, j int) bool {
		return Freqs[i].Freq > Freqs[j].Freq
	})

	for i := 0; i < TOP_N && i < len(Freqs)-1; i++ {
		fmt.Println(Freqs[i])
	}

}

// extractNgrams tokenizes the supplied stream and
// extracts each unique n-gram into a map relating
// the n-gram to its frequency
func extractNgrams(r io.Reader) {

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	wordQueue := make([]string, GRAMSIZE)

	for scanner.Scan() {
		word := scrubWord(scanner.Text())
		wordQueue = append(wordQueue, word)
		if len(wordQueue) > GRAMSIZE {
			wordQueue = wordQueue[1:] // Drop the first element
		}
		if len(wordQueue) == GRAMSIZE {
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

var invalidChars = regexp.MustCompile("[^a-zA-Z0-9]+")

// scrubWord scrubs noise characters (punctuation, etc)
// and lowercases the input
func scrubWord(s string) string {

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
