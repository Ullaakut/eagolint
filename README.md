# Eagolint

<p align="center">
    <a href="http://img.shields.io/badge/godoc-reference-blue.svg?style=flat">
        <img src="https://godoc.org/github.com/Ullaakut/eagolint"/>
    </a>
    <a href="http://img.shields.io/badge/license-MIT-red.svg?style=flat">
        <img src="https://raw.githubusercontent.com/Ullaakut/eagolint/master/LICENSE"/>
    </a>
    <a href="https://goreportcard.com/badge/github.com/Ullaakut/eagolint">
        <img src="https://goreportcard.com/report/Ullaakut/eagolint"/>
    </a>
</p>

Eagolint is a go linter that keeps your comments punctuated and double-space free.

## Usage

```text
A go linter that checks for comments with missing punctuation and double spaces.

Usage:
  eagolint [flags] [path ...]

Flags:
  -e, --exclude string      Exclude lines that match this regex
      --files               Read file names from stdin
  -g, --go-only             Only check .go files
  -h, --help                help for eagolint
  -s, --skip-list strings   List of directories to skip
  -t, --skip-tests          Skip _test.go files
      --vendor              Check files in vendor directory
```

## License

This project is under the MIT license. See [LICENSE](LICENSE) for more information.
