package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	SUCCEED int = iota
	NOMATCH
	FAILURE
)

// Build metadata
var (
	Version = "0.0"
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
)

type options struct {
	recursive  bool
	count      bool
	lineNumber bool
	ignoreCase bool
}

func run(opts options, pattern string, files []string) {
	/* Later implementation */
}

// Remember `main()` exited with a 0 exit code
func main() {

	// Must be called before flag were accessed by the program
	flag.Parse()

	// Override Usage() with our own layout
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n  grab [options] <pattern> [file1 file2 ...]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\nOptions:\n")
		flag.PrintDefaults()
	}

	if *showVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "grab %s\ncommit: %s\nbuilt: %s\n", Version, Commit, Date)
		os.Exit(SUCCEED)
	}

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(FAILURE)
	}

	pattern := flag.Arg(0)
	files := flag.Args()[1:]

	opts := options{
		recursive:  *recursive,
		count:      *count,
		lineNumber: *lineNumber,
		ignoreCase: *ignoreCase,
	}

	run(opts, pattern, files)
}
