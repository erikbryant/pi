package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	find = flag.String("find", "pi", "string to search")
)

const piDir = "../datafiles/pi/"

// readPackedPi100 returns a byte slice of digits of pi
func readPackedPi100() ([]byte, error) {
	fileName := "100million.bin"
	file := path.Join(piDir, fileName)

	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// stringToBytes returns a byte slice representing the ASCII value of the string
func stringToBytes(s string) []byte {
	if len(s) == 0 {
		panic("Empty string!")
	}

	// Convert the characters in str to their ASCII base 10 digits
	digits := ""
	for _, ch := range s {
		digits += fmt.Sprintf("%d", ch)
	}

	// Convert the digits to their numeric value
	bytes := []byte{}
	for _, ch := range digits {
		bytes = append(bytes, byte(ch)-'0')
	}

	return bytes
}

// digit returns the nth digit of the packed bytes
func digit(pi []byte, i int) byte {
	offset := i >> 1

	if offset >= len(pi) {
		panic("Beyond end of Pi's digits")
	}

	if i&0x01 == 0 {
		// An even position
		return (pi[offset] & 0xf0) >> 4
	}

	// An odd position
	return pi[offset] & 0x0f
}

// unpack unpacks the given byte slice and returns it
func unpack(pi []byte, start, length int) []byte {
	unpacked := []byte{}

	for i := start; i < start+length; i++ {
		unpacked = append(unpacked, digit(pi, i))
	}

	return unpacked
}

// searchPi searches the digits of pi for a match to target
func searchPi(pi, target []byte) {
	// For every digit in Pi...
	for i := 0; i < len(pi)*2; i++ {
		if digit(pi, i) == target[0] {
			// We have the start of a match!
			found := true
			for j := 1; j < len(target); j++ {
				if digit(pi, i+j) != target[j] {
					found = false
					break
				}
			}
			if found {
				fmt.Printf("Found a match! Target: %v, Digit: %d, Context: %v\n", target, i, unpack(pi, i, len(target)))
			}
		}
	}
}

// search searches the digits of pi for a match for s
func search(s string) {
	target := stringToBytes(s)
	fmt.Printf("Searching PI for: %s -> %v\n", s, target)

	pi, err := readPackedPi100()
	if err != nil {
		panic(err)
	}

	searchPi(pi, target)
}

func main() {
	fmt.Printf("Welcome to String Finder!\n\n")

	flag.Parse()

	search(*find)
	search(strings.ToUpper(*find))
}
