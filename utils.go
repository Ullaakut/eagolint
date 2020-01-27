package eagolint

import (
	"bufio"
	"bytes"
	"strings"
)

var (
	genHdr = []byte("// Code generated ")
	genFtr = []byte(" DO NOT EDIT.")
)

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
func isGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, genHdr) && bytes.Contains(b, genFtr) {
			return true
		}
	}
	return false
}

// contains returns whether the given string slice contains the given string.
func contains(slice []string, str string) bool {
	for _, sliceElem := range slice {
		if str == sliceElem {
			return true
		}
	}
	return false
}

// Punctuated returns whether a line ends with punctuation or a
// closing parentheses/bracket.
func isPunctuated(line string) bool {
	return strings.HasSuffix(line, ".") ||
		strings.HasSuffix(line, "!") ||
		strings.HasSuffix(line, "?") ||
		strings.HasSuffix(line, "}") ||
		strings.HasSuffix(line, "]") ||
		strings.HasSuffix(line, ")")
}
