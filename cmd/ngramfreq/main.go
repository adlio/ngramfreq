package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

// GramSize is the target length of word
// sequences (i.e. 3 means we search for
// trigrams)
var GramSize = 3

// OutputSize is the number of n-grams to
// output
var OutputSize = 100

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

	scanner := NewScanner(GramSize)

	for _, filename := range filenames {
		scanner.ScanFile(filename)
	}

	if haveStdIn {
		scanner.Scan(os.Stdin)
	}

	scanner.Score()
	scanner.WriteTopN(OutputSize, os.Stdout)
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
