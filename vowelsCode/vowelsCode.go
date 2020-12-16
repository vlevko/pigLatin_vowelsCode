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
)

var matches = map[string]string{
	"a": "1",
	"e": "2",
	"i": "3",
	"o": "4",
	"u": "5",
	"A": "1",
	"E": "2",
	"I": "3",
	"O": "4",
	"U": "5",
	"1": "a",
	"2": "e",
	"3": "i",
	"4": "o",
	"5": "u",
}

var decodeString func(string, func(string) string) string
var encodeString func(string, func(string) string) string

// init initializes regexp patterns with appropriate functions
func init() {
	decodeString = regexp.MustCompile(`([1-5]){1}`).ReplaceAllStringFunc
	encodeString = regexp.MustCompile(`([aeiouAEIOU]){1}`).ReplaceAllStringFunc
}

// main runs decode/encode utility, parses a flag, gets a phrase, decodes/
// encodes it, prints result
func main() {
	decode := flag.Bool("d", false, "decode phrase mode; default is encode mode")
	flag.Parse()
	phrase := getPhraseFromStdin()

	if *decode {
		phrase = decodeString(phrase, replace)
	} else {
		phrase = encodeString(phrase, replace)
	}

	fmt.Print(phrase)
}

// getPhraseFromStdin handles user input from os.Stdin, returns it as a string,
// if error, logs it and exits with non-zero code
func getPhraseFromStdin() string {
	fmt.Println("Enter a phrase:")
	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')

	if err != nil && err != io.EOF {
		log.Fatalf("piglatin: error: %v\n", err)
	}

	return phrase
}

// replace looks for matching value by key to replace with, gets a string,
// returns the match as a string
func replace(s string) string {
	return matches[s]
}
