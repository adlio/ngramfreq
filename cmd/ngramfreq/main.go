package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

const GRAMSIZE = 3

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
}

// extractNgrams tokenizes the supplied stream and
// extracts each unique n-gram into a map relating
// the n-gram to its frequency
func extractNgrams(r io.Reader) map[string]int {

	grams := make(map[string]int)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	wordQueue := make([]string, GRAMSIZE)

	for scanner.Scan() {
		word := scanner.Text()
		wordQueue = append(wordQueue, word)
		if len(wordQueue) > GRAMSIZE {
			wordQueue = wordQueue[1:] // Drop the first element
		}
		if len(wordQueue) == GRAMSIZE {
			phrase := strings.Join(wordQueue, " ")
			grams[phrase]++
		}
	}

	return grams
}
