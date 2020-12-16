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
	"unicode/utf8"
)

const vowels = "aeiouAEIOU"

var getAllMatches func(string, int) []string
var isEnglishLetter func(string) bool

// init initializes regexp patterns with appropriate functions
func init() {
	getAllMatches = regexp.MustCompile(`([a-zA-Z]+)|([^a-zA-Z]*)`).FindAllString
	isEnglishLetter = regexp.MustCompile(`[a-zA-Z]`).MatchString
}

// main runs Pig Latin utility, gets a phrase, translates it, prints translation
func main() {
	phrase := getPhraseFromStdin()
	phrase = translatePhrase(phrase)
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

// translatePhrase gets a string, translates it concurrently, returns
// translation as a string
func translatePhrase(phrase string) string {
	// matches words and symbols in phrase
	matches := getAllMatches(phrase, -1)
	// chooses a suitable number of goroutines
	goroutinesCount := minInt(len(matches), runtime.NumCPU())
	// uses buffered channel for the limited number of concurrent goroutines
	ch := make(chan *string, goroutinesCount)
	// uses wait group for completion of all gourutines
	var wg sync.WaitGroup
	// loops through the number of goroutines and spawns translate goroutines
	for i := 0; i < goroutinesCount; i++ {
		go translateGoroutine(ch, &wg)
	}
	// loops through all matched parts and sends them to translate goroutines
	for i := range matches {
		wg.Add(1)
		ch <- &matches[i]
	}
	// waits for goroutines complete translation and terminates theme
	wg.Wait()
	close(ch)
	// joins all translated words with symbols and returns translation
	return strings.Join(matches, "")
}

// minInt gets two integers, returns the minimum one as an integer
func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// translateGoroutine handles requests for translation, gets a channel of
// pointer to string to read from and a pointer to wait group for notification
// of completion
func translateGoroutine(ch chan *string, wg *sync.WaitGroup) {
	for match := range ch {
		r, size := utf8.DecodeRuneInString(*match)

		if isEnglishLetter(string(r)) {
			*match = translateWord(*match, r, size)
		}

		(*wg).Done()
	}
}

// translateWord gets a string for translation with the first rune and its size,
// translates it, returns translation as a string
func translateWord(word string, r rune, size int) string {
	// appends a suffix if word starts with vowel and returns
	if strings.ContainsRune(vowels, r) {
		return word + "yay"
	}
	// slices consonants from start to end
	consonants := word[:size]
	word = word[size:]
	// loops through all consonants in a row
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		// stops slicing if vowel found
		if strings.ContainsRune(vowels, r) {
			break
		}
		// slices one rune forward
		consonants += word[:size]
		word = word[size:]
	}
	// appends a suffix and return
	return word + consonants + "ay"
}
