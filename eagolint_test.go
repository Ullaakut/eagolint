package eagolint_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/Ullaakut/eagolint"
	"io/ioutil"
)

func TestShouldSkipDirs(t *testing.T) {
	skip, err := eagolint.ShouldSkip(".git", true, []string{".git"}, false, false)
	if skip == false || err != filepath.SkipDir {
		t.Errorf("Expected %t, %s got. %t, %s", true, filepath.SkipDir, skip, err)
	}

	skip, err = eagolint.ShouldSkip("dir", true, []string{".git"}, false, false)
	if skip == false || err != nil {
		t.Errorf("Expected %t, %v got. %t, %s", true, nil, skip, err)
	}
}

func TestShouldSkipFiles(t *testing.T) {
	t.Run("regular files", func(t *testing.T) {
		binaryFilePath, _ := os.Executable()
		tests := []struct {
			path      string
			goOnly    bool
			skipTests bool

			shouldSkip bool
			err        error
		}{
			{path: "eagolint.go", goOnly: false, skipTests: false, shouldSkip: false},
			{path: "eagolint.go", goOnly: true, skipTests: false, shouldSkip: false},
			{path: "eagolint.go", goOnly: false, skipTests: true, shouldSkip: false},
			{path: "eagolint.go", goOnly: true, skipTests: true, shouldSkip: false},
			{path: "README.md", goOnly: false, skipTests: false, shouldSkip: false},
			{path: "README.md", goOnly: true, skipTests: false, shouldSkip: true},
			{path: "README.md", goOnly: false, skipTests: true, shouldSkip: false},
			{path: "README.md", goOnly: true, skipTests: true, shouldSkip: true},
			{path: "eagolint_test.go", goOnly: false, skipTests: false, shouldSkip: false},
			{path: "eagolint_test.go", goOnly: true, skipTests: false, shouldSkip: false},
			{path: "eagolint_test.go", goOnly: false, skipTests: true, shouldSkip: true},
			{path: "eagolint_test.go", goOnly: true, skipTests: true, shouldSkip: true},
			{path: binaryFilePath, goOnly: false, skipTests: false, shouldSkip: true},
			{path: binaryFilePath, goOnly: true, skipTests: false, shouldSkip: true},
			{path: binaryFilePath, goOnly: false, skipTests: true, shouldSkip: true},
			{path: binaryFilePath, goOnly: true, skipTests: true, shouldSkip: true},
		}

		for i, tc := range tests {
			skip, err := eagolint.ShouldSkip(tc.path, false, []string{".git"}, tc.goOnly, tc.skipTests)
			if skip != tc.shouldSkip || err != tc.err {
				t.Errorf("%d) Expected %t, %v got %t, %s", i+1, tc.shouldSkip, tc.err, skip, err)
			}
		}
	})
	t.Run("file in skiplist", func(t *testing.T) {
		skip, err := eagolint.ShouldSkip("file", false, []string{"file"}, false, false)
		if skip != true || err != nil {
			t.Errorf("Expected %t, %v got. %t, %s", true, nil, skip, err)
		}
	})
	t.Run("error on file not found", func(t *testing.T) {
		skip, err := eagolint.ShouldSkip("file", false, []string{".git"}, false, false)
		if skip != true || err == nil {
			t.Errorf("Expected %t, %v got. %t, %s", true, nil, skip, err)
		}
	})
}

func TestProcess(t *testing.T) {
	content, err := ioutil.ReadFile("test_assets/bad_comments.go")
	if err != nil {
		t.Fatal(err)
	}

	b := &bytes.Buffer{}
	err = eagolint.Process(bytes.NewBuffer(content), b, "file", nil)
	if err != nil {
		t.Errorf("Expected %v, got %s", nil, err)
	}

	expected := `file:6: double space typo in comment
file:8: double space typo in comment
file:12: missing punctuation at end of comment
file:14: missing punctuation at end of comment
file:20: double space typo in comment
file:21: double space typo in comment
file:25: missing punctuation at end of comment
file:31: double space typo in comment
`
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func TestProcessFile(t *testing.T) {
	b := &bytes.Buffer{}
	err := eagolint.ProcessFile(b, "eagolint_test.go", nil)
	if err != nil {
		t.Errorf("Expected %v, got %s", nil, err)
	}
}

func TestProcessExclude(t *testing.T) {
	lines := `	// TODO: fix
				// FIXME: do something
				// This is a non-excluded  comment with issues`
	b := &bytes.Buffer{}
	exclude := regexp.MustCompile("TODO|FIXME")
	expected := `file:3: missing punctuation at end of comment
file:3: double space typo in comment
`
	_ = eagolint.Process(bytes.NewBufferString(lines), b, "file", exclude)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}
