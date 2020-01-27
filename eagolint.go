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

type line struct {
	pos int
	text string
}

// Process checks all lines in the reader and writes an error if the line length
// is greater than MaxLength.
func Process(r io.Reader, w io.Writer, path string, exclude *regexp.Regexp) error {
	var pos int
	s := bufio.NewScanner(r)

	var comments []line
	for s.Scan() {
		pos++
		text := strings.TrimSpace(s.Text())

		isComment := strings.HasPrefix(text, "//")
		if !isComment {
			continue
		}

		comments = append(comments, line{
			pos: pos,
			text: text,
		})
	}

	processComments(w, comments, path)

	return s.Err()
}

// processComments processes all comments from a file.
func processComments(w io.Writer, c []line, path string) {
	// First, split the map of all comment lines into clusters.
	// (1 cluster = 1 multiline/inline comment)
	var (
		clusters [][]line
		currentCluster []line
		prev = -1
	)
	for _, l := range c {
		// This is a new cluster, so we append the previous cluster to the list
		// of comment clusters.
		if l.pos != prev+1 {
			clusters = append(clusters, currentCluster)
			currentCluster = []line{}
		}

		currentCluster = append(currentCluster, l)
		prev = l.pos
	}

	if len(currentCluster) > 0 {
		clusters = append(clusters, currentCluster)
	}

	// Process each comment cluster individually.
	for _, cluster := range clusters {
		processComment(w, cluster, path)
	}
}

// processComment processes an inline comment or multiline comment block.
func processComment(w io.Writer, comment []line, path string) {
	// Iterate on each time that is part of the comment.
	for idx, line := range comment {
		// If this is the last line of the comment and it's missing punctuation,
		// print a warning.
		if idx+1 == len(comment) {
			if !isPunctuated(line.text) {
				_, _ = fmt.Fprintf(w, "%s:%d: missing punctuation at end of comment\n", path, line.pos)
			}
		}

		// If the line contains a double space, print a warning.
		if strings.Contains(line.text, "  ") {
			_, _ = fmt.Fprintf(w, "%s:%d: double space typo in comment\n", path, line.pos)
		}
	}
}
