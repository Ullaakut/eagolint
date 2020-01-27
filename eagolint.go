package eagolint

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ShouldSkip checks a path and determines if it should be skipped.
// SkipList contains paths to be skipped skip.
// All directories are skipped, only files are processed.
// If GoOnly is true, non-go files are skipped.
// Otherwise, checks that file is readable text file.
func ShouldSkip(path string, isDir bool, skipList []string, goOnly bool, skipTests bool) (bool, error) {
	name := filepath.Base(path)
	if contains(skipList, name) {
		if isDir {
			return true, filepath.SkipDir
		}
		return true, nil
	}

	if isDir {
		return true, nil
	}

	if skipTests && strings.HasSuffix(path, "_test.go") {
		return true, nil
	}

	isGo := strings.HasSuffix(path, ".go")
	if goOnly && !isGo {
		return true, nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return true, err
	}

	if isGo {
		return isGenerated(b), nil
	}

	m := http.DetectContentType(b)
	if !strings.Contains(m, "text/") {
		return true, nil
	}

	return false, nil
}

// ProcessFile checks all lines in the file and writes an error if the line
// length is greater than MaxLength.
func ProcessFile(w io.Writer, path string, exclude *regexp.Regexp) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error closing file: %s\n", err)
		}
	}()

	return Process(f, w, path, exclude)
}

// Process checks all lines in the reader and writes an error if the line length
// is greater than MaxLength.
func Process(r io.Reader, w io.Writer, path string, exclude *regexp.Regexp) error {
	l := 0
	s := bufio.NewScanner(r)

	var (
		// prevIsComment is used to track whether the previous line is a comment,
		// in order to detect multi-line comment blocks.
		prevIsComment bool

		// punctuated stores whether the last comment line found was punctuated.
		// if a comment ends without punctuation, a warning is printed.
		punctuated bool
	)
	for s.Scan() {
		l++
		line := strings.TrimSpace(s.Text())
		isComment := strings.HasPrefix(line, "//")

		if !isComment {
			if prevIsComment && !punctuated {
				// Previous line was a comment and it ended on this line.
				// If punctuated is false, then a warning is printed for the
				// previous line.
				_, _ = fmt.Fprintf(w, "%s:%d: missing punctuation at end of comment\n", path, l-1)
			}

			prevIsComment = false
			continue
		}

		if strings.Contains(line, "  ") {
			_, _ = fmt.Fprintf(w, "%s:%d: double space typo in comment\n", path, l)
		}

		punctuated = isPunctuated(line)
		prevIsComment = true
	}

	return s.Err()
}
