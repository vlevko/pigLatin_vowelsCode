// Package main implements a decode/encode utility for vowels
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var matches = map[string]string{
	"a": "1",
	"e": "2",
	"i": "3",
	"o": "4",
	"u": "5",
	"1": "a",
	"2": "e",
	"3": "i",
	"4": "o",
	"5": "u",
}

var reDecode *regexp.Regexp
var reEncode *regexp.Regexp

// Initialize regexp patterns
func init() {
	reDecode = regexp.MustCompile(`([1-5]){1}`)
	reEncode = regexp.MustCompile(`([aeiouAEIOU]){1}`)
}

// Run decode/encode utility
func main() {
	// Parse flag
	decode := flag.Bool("d", false, "decode phrase mode; default is encode mode")
	flag.Parse()
	// Get phrase
	phrase := getPhraseFromStdin()
	// Decode/encode phrase
	if *decode {
		phrase = reDecode.ReplaceAllStringFunc(phrase, repl)
	} else {
		phrase = reEncode.ReplaceAllStringFunc(phrase, repl)
	}
	// Print encoded/decoded phrase
	fmt.Print(phrase)
}

// Handle user input from os.Stdin, return it as a string
func getPhraseFromStdin() string {
	fmt.Println("Enter a phrase:")
	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')
	// If error, log it and exit with non-zero code
	if err != nil && err != io.EOF {
		log.Fatalf("piglatin: error: %v\n", err)
		os.Exit(1)
	}

	return phrase
}

// Look for matching key:value to replace with
func repl(s string) string {
	return matches[strings.ToLower(s)]
}
