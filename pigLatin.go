// Package main implements a translate utility from English to Pig Latin
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

const vowels = "aeiouAEIOU"

var re *regexp.Regexp

// Initialize regexp pattern
func init() {
	re = regexp.MustCompile(`(([a-zA-Z]+)|([^a-zA-Z]*))`)
}

// Run Pig Latin utility
func main() {
	// Get phrase
	phrase := getPhraseFromStdin()
	// Translate phrase
	phrase = translatePhrase(phrase)
	// Print translated phrase
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

// Translate a phrase concurrently, return translation as a string
func translatePhrase(phrase string) string {
	// Match words and symbols in phrase
	matches := re.FindAllString(phrase, -1)
	// Choose a suitable number of goroutines
	goroutinesCount := min(len(matches), runtime.NumCPU())
	// Use buffered channel to limit the number of concurrent goroutines
	ch := make(chan *string, goroutinesCount)
	// Use wait group for completion of all gourutines
	var wg sync.WaitGroup
	// Loop through the number of goroutines and spawn translate goroutines
	for i := 0; i < goroutinesCount; i++ {
		go translateGoroutine(ch, &wg)
	}
	// Loop through all matched parts and send them to translate goroutines
	for i := range matches {
		wg.Add(1)
		ch <- &matches[i]
	}
	// Wait for goroutines complete translation and terminate theme
	wg.Wait()
	close(ch)
	// Join all translated words with symbols and return
	return strings.Join(matches, "")
}

// Return minimum of two integer numbers
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Handle concurrent requests for translation
func translateGoroutine(ch chan *string, wg *sync.WaitGroup) {
	for match := range ch {
		r, size := utf8.DecodeRuneInString(*match)
		if unicode.IsLetter(r) {
			*match = translateWord(*match, r, size)
		}
		(*wg).Done()
	}
}

// Translate a word, return it as a string
func translateWord(word string, r rune, size int) string {
	// Append a suffix if word starts with vowel and return
	if strings.ContainsRune(vowels, r) {
		return word + "yay"
	}
	// Slice consonants from start to end
	consonants := word[:size]
	word = word[size:]
	// Loop through all consonants in a row
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		// Stop slicing if vowel found
		if strings.ContainsRune(vowels, r) {
			break
		}
		// Slice one rune forward
		consonants += word[:size]
		word = word[size:]
	}
	// Append a suffix and return
	return word + consonants + "ay"
}
