package main

import (
	"fmt"
	"os"

	"bufio"
	"github.com/Ullaakut/eagolint"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	"regexp"
)

var args struct {
	GoOnly    bool     `arg:"-g,env,help:only check .go files"`
	SkipTests bool     `arg:"-t,env,help:skip _test.go files"`
	Input     []string `arg:"positional"`
	SkipList  []string `arg:"-s,env,help:list of dirs to skip"`
	Vendor    bool     `arg:"env,help:check files in vendor directory"`
	Files     bool     `arg:"help:read file names from stdin one at each line"`
	Exclude   string   `arg:"-e,env,help:exclude lines that matches this regex"`
}

var cmd = &cobra.Command{
	Use:   "eagolint",
	Short: "Eagolint is a go linter that keeps your comments punctuated and double-space free.",
	Long:  `A go linter that checks for comments with missing punctuation and double spaces.`,
	Run: func(cmd *cobra.Command, input []string) {
		args.GoOnly = viper.GetBool("go-only")
		args.SkipTests = viper.GetBool("skip-tests")
		args.Vendor = viper.GetBool("vendor")
		args.Files = viper.GetBool("files")
		args.SkipList = viper.GetStringSlice("skip-list")
		args.Exclude = viper.GetString("exclude")
		args.Input = input

		run()
	},
}

func init() {
	viper.SetEnvPrefix("eagolint")
	viper.AutomaticEnv()

	cmd.Flags().BoolP("go-only", "g", false, "Only check .go files")
	cmd.Flags().BoolP("skip-tests", "t", false, "Skip _test.go files")
	cmd.Flags().Bool("vendor", false, "Check files in vendor directory")
	cmd.Flags().Bool("files", false, "Read file names from stdin")
	cmd.Flags().StringSliceP("skip-list", "s", nil, "List of directories to skip")
	cmd.Flags().StringP("exclude", "e", "", "Exclude lines that match this regex")

	_ = viper.BindPFlags(cmd.Flags())
}

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() {
	// Ensure that the provided exclusion regexp compiles.
	var exclude *regexp.Regexp
	if args.Exclude != "" {
		e, err := regexp.Compile(args.Exclude)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error compiling exclude regexp: %s\n", err)
			os.Exit(1)
		}
		exclude = e
	}

	// Remove the vendor dir from the skip list if --vendor is true.
	if args.Vendor {
		for i, p := range args.SkipList {
			if p == "vendor" {
				args.SkipList = append(args.SkipList[:i], args.SkipList[:i]...)
			}
		}
	}

	// If --files is set to true, process each line from stdin as a file and exit.
	if args.Files {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			err := eagolint.ProcessFile(os.Stdout, s.Text(), exclude)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error processing file: %s\n", err)
			}
		}
		os.Exit(0)
	}

	// Otherwise, walk the inputs recursively.
	for _, d := range args.Input {
		err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
			if info == nil {
				_, _ = fmt.Fprintf(os.Stderr, "eagolint: %s no such file or directory\n", path)
				return nil
			}
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "eagolint: %s\n", err)
				return nil
			}
			skip, ret := eagolint.ShouldSkip(path, info.IsDir(), args.SkipList, args.GoOnly, args.SkipTests)
			if skip {
				return ret
			}

			return eagolint.ProcessFile(os.Stdout, path, exclude)
		})

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error walking the file system: %s\n", err)
			os.Exit(1)
		}
	}
}
