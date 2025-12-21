package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

const (
	SUCCEED int = iota
	NOMATCH
	FAILURE
)

// Build metadata
var (
	Version = "v0.0"
	Commit  = "none"
	Date    = "unknown"
)

var showVersion = flag.Bool("version", false, "Show version information")

// Behavioral flags: affects how the program ran
var (
	recursive  = flag.Bool("r", false, "Recursively search directories")
	count      = flag.Bool("c", false, "Show count of matching lines")
	lineNumber = flag.Bool("n", false, "Show line numbers")
	ignoreCase = flag.Bool("i", false, "Ignore case distinction")
	color      = flag.String("color", "auto", "Show colors on matching lines, is either 'always', 'auto', 'never'")
)

type options struct {
	recursive  bool
	count      bool
	lineNumber bool
	ignoreCase bool
	color      string
}

// Remember `main.main` normally exited with a zero exit code
func main() {

	// Use the default FlagSet 'CommandLine'
	flag.Parse()

	// Override Usage() with our own
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n  grab [options] <pattern> [file1 file2 ...]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
		flag.PrintDefaults()
	}

	if *showVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "grab %s\ncommit: %s\nbuilt: %s\n", Version, Commit, Date)
		os.Exit(SUCCEED)
	}

	// If no pattern were given exit
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(FAILURE)
	}

	pattern := flag.Arg(0)
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "grab %s: invalid pattern\n", pattern)
		os.Exit(FAILURE)
	}

	files := flag.Args()[1:]

	opts := options{
		recursive:  *recursive,
		count:      *count,
		lineNumber: *lineNumber,
		ignoreCase: *ignoreCase,
		color:      *color,
	}
	_ = opts

	// When no input files are provided, read from stdin
	// This enables usage such as piping `cat file | grab pattern`
	if len(files) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		matched := false

		for scanner.Scan() {
			line := scanner.Bytes()

			// Skip lines that has no matches
			if !re.Match(line) {
				continue
			}

			matched = true

			// Replace all matches with colored equivalents
			coloredLine := re.ReplaceAllFunc(line, func(b []byte) []byte {
				return append(
					append([]byte("\033[1;31m"), b...),
					"\033[0m"...,
				)
			})

			// Write directly to stdout
			os.Stdout.Write(coloredLine)
			os.Stdout.Write([]byte{'\n'})
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "grab: %v\n", err)
			os.Exit(FAILURE)
		}

		// If no matches were found in stdin, exit with a non-zero status
		// This behavior is consistent with the way unix grep works
		if !matched {
			os.Exit(NOMATCH)
		}
	}
}
