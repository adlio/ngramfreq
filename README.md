> Stackery.io Coding Exercise - adlio
# ngramfreq

This program calculates the most frequently occuring [n-grams](https://en.wikipedia.org/wiki/N-gram) in text files supplied as command line arguments or via standard input. By default it looks for *3*-word sequences and outputs the *100* most common found.

### Usage Examples

Command line arguments:

		ngramfreq file1.txt file2.txt

Standard input and command line args combined:

		cat first.txt | ngramfreq second.txt

Search for bigrams instead of trigrams:

		ngramfreq -s 2 file1.txt

Output the top 10 instead of the top 100:

		ngramfreq -n 10 file1.txt

## Developer Setup Guide

This application is written in Go, using only the Go standard library. Follow the [Go install guide](https://golang.org/doc/install#install) appropriate for your OS.

**Prerequisites**:

* Go (1.8+)

### Cloning the Project

You'll want to check the project out in your `$GOPATH`.

		mkdir -p $GOPATH/src/github.com/samg
		git clone git@github.com:samg/adlio-coding-exercise
		cd adlio-coding-exercise

### Running Tests

Tests should be run from the `cmd/ngramfreq` directory:

			cd cmd/ngramfreq
			go test

### Compiling and Running

Compile and run the app from the `cmd/ngramfreq` directory:

			cd cmd/ngramfreq
			go build
			./ngramfreq file1.txt

### Architecture Decision Record (ADR)

This project's architectural decisions are being tracked in Markdown files in the [adr](./adr) directory.

----

## Application Requirements

Below are the original requirements for this project:

1. The program can be written in any language but should be easy to run.
2.It should accept as arguments a list of one or more file paths (e.g. ./solution file1.txt file2.txt ...). The program also accepts input on stdin (e.g. cat file1.txt | ./solution).
3. The program outputs a list of the 100 most common three word sequences in the text, along with a count of how many times each occurred in the text.

	For example:
	
		231 - i will not
		116 - i do not
		105 - there is no
		54 - i know not
		37 - i am not

4. The program ignores punctuation, line endings, and is case insensitive (e.g. `I love\nsandwiches.` should be treated the same as `(I LOVE SANDWICHES!!)`).
5. The program is capable of processing large files and runs as fast as possible.
