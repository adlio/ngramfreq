package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
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
	haveStdin, filenames, err := processArgs()
	if err != nil {
		usageErr(err, os.Stderr)
		os.Exit(1)
	}

	scanner := NewScanner(GramSize)
	for _, filename := range filenames {
		scanner.ScanFile(filename)
	}
	if haveStdin {
		scanner.Scan(os.Stdin)
	}
	scanner.Score()
	scanner.WriteTopN(OutputSize, os.Stdout)
}

func processArgs() (haveStdin bool, filenames []string, err error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.IntVar(&GramSize, "s", 3, "sequence size (2 = bigrams, 3=trigrams)")
	fs.IntVar(&OutputSize, "n", 100, "number of n-grams to output")
	err = fs.Parse(os.Args[1:])
	if err != nil {
		return
	}

	// Get a list of valid files
	var info os.FileInfo
	for _, arg := range fs.Args() {
		filename := path.Clean(arg)
		if info, err = os.Stat(filename); err == nil && info.Mode().IsRegular() {
			filenames = append(filenames, filename)
		} else {
			err = fmt.Errorf("Can't process '%s' as a file", filename)
			return
		}
	}

	// Determine if we have character input on os.Stdin
	if fi, _ := os.Stdin.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
		haveStdin = true
	}

	if len(filenames) == 0 && !haveStdin {
		err = fmt.Errorf("Either a filename or standard input is required")
	}

	return
}

func usageErr(err error, w io.Writer) {
	fmt.Fprintf(w, "%s\n\n", err)
	fmt.Fprintf(w, "Usage: ngramfreq [option flags] file1.txt file2.txt ...\n")
	fmt.Fprintf(w, "-s int  sequence size (2=bigrams, 3=trigrams) (default: 3)\n")
	fmt.Fprintf(w, "-n int  output size (default: 100)\n")
}
