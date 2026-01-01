package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aelberthcheong/grab/internal/fileutil"
	"github.com/aelberthcheong/grab/internal/matcher"
)

// Exit codes
const (
	Succeed int = iota
	NoMatch
	Failure
)

// Metadata
var (
	Version = "v0.0"
	Date    = "unknown"
)

const License = `
MIT License

Copyright (c) 2025 Aelberth Cheong

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiBlue   = "\033[34m"
	ansiPurple = "\033[35m"
)

// Flags
var showVersion = flag.Bool("version", false, "Show version information")

var (
	recursive    = flag.Bool("r", false, "Recursively search directories")
	count        = flag.Bool("c", false, "Show count of matching lines")
	lineNumber   = flag.Bool("n", false, "Show line numbers")
	ignoreCase   = flag.Bool("i", false, "Ignore case distinction")
	withFilename = flag.Bool("h", false, "Show file name with output lines")
	invertMatch  = flag.Bool("v", false, "Select non-matching lines")
	noErrMessage = flag.Bool("s", false, "Suppress error messages")
	noColorLine  = flag.Bool("color", false, "Show colors on matching lines")
)

type options struct {
	recursive    bool
	count        bool
	lineNumber   bool
	ignoreCase   bool
	withFilename bool
	invertMatch  bool
	noErrMessage bool
	noColorLine  bool
}

type Emitter struct {
	opts    options
	re      *regexp.Regexp
	matches int
}

func (e *Emitter) Emit(line []byte, lineNo int, filename string) {
	e.matches++

	if e.opts.count {
		return
	}

	if e.opts.noColorLine {
		// No color output
		if e.opts.withFilename {
			fmt.Printf("%s:", filename)
		}
		if e.opts.lineNumber {
			fmt.Printf("%d:", lineNo)
		}
		fmt.Println(string(line))
		return
	}

	// Colored output
	if e.opts.withFilename {
		fmt.Printf("%s%s%s", ansiPurple, filename, ansiReset)
		fmt.Printf("%s:%s", ansiBlue, ansiReset)
	}

	if e.opts.lineNumber {
		fmt.Printf("%s%d%s", ansiGreen, lineNo, ansiReset)
		fmt.Printf("%s:%s", ansiBlue, ansiReset)
	}

	colored := highlight(line, e.re)
	fmt.Println(string(colored))
}

func (e *Emitter) Finalize() {
	if e.opts.count {
		fmt.Println(e.matches)
	}
}

func walkDir(root string, m *matcher.Matcher, emitter *Emitter, opts options) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if !opts.noErrMessage {
				fmt.Fprintf(os.Stderr, "grab: %v\n", err)
			}
			return nil // continue traversal
		}

		// Skip directories; WalkDir handles recursion
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			if !opts.noErrMessage {
				fmt.Fprintf(os.Stderr, "grab: %v\n", err)
			}
			return nil
		}
		defer f.Close()

		if !fileutil.IsLikelyText(f) {
			return nil
		}
		_, _ = f.Seek(0, io.SeekStart)

		_, _ = processReader(f, path, m, emitter)
		return nil
	})
}

func highlight(line []byte, re *regexp.Regexp) []byte {
	indexes := re.FindAllIndex(line, -1)
	if len(indexes) == 0 {
		return line
	}

	var out []byte
	last := 0

	for _, idx := range indexes {
		start, end := idx[0], idx[1]
		out = append(out, line[last:start]...)
		out = append(out, ansiRed...)
		out = append(out, line[start:end]...)
		out = append(out, ansiReset...)
		last = end
	}

	out = append(out, line[last:]...)
	return out
}

func processReader(r io.Reader, filename string, matcher *matcher.Matcher, emitter *Emitter) (bool, error) {
	scanner := bufio.NewScanner(r)

	buf := make([]byte, 0, bufio.MaxScanTokenSize) // By default a token is a line
	scanner.Buffer(buf, 1024*1024)                 // Increase buffer size from ~65KB to ~1MB

	var (
		lineNum int
		found   bool
	)

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if !matcher.Match(line) {
			continue
		}

		found = true

		// Cause we reuse `Scanner.Bytes()` we need to copy the line
		copiedLine := make([]byte, len(line))
		copy(copiedLine, line)

		emitter.Emit(copiedLine, lineNum, filename)
	}

	return found, scanner.Err()
}

// Remember `main.main` normally exited with a zero exit code
func main() {

	// Override Usage() with our own
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  grab [options] <pattern> [file1 file2 ...]\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	// Use the default FlagSet 'CommandLine'
	flag.Parse()

	if *showVersion {
		fmt.Printf("grab %s\nAuthor: Aelberth Cheong <aelberth.cheong@outlook.com>\nBuilt : %s\n", Version, Date)
		fmt.Printf(License)
		os.Exit(Succeed)
	}

	// If no pattern were given exit
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(Failure)
	}

	opts := options{
		recursive:    *recursive,
		count:        *count,
		lineNumber:   *lineNumber,
		ignoreCase:   *ignoreCase,
		withFilename: *withFilename,
		invertMatch:  *invertMatch,
		noErrMessage: *noErrMessage,
		noColorLine:  *noColorLine,
	}

	pattern := flag.Arg(0)
	mtch, err := matcher.NewMatcher(pattern, opts.ignoreCase, opts.invertMatch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "grab: invalid pattern\n")
		os.Exit(Failure)
	}

	emitter := &Emitter{opts: opts, re: mtch.Regexp()}

	files := flag.Args()[1:]
	anyFound := false

	if len(files) == 0 {
		found, err := processReader(os.Stdin, "(standard input)", mtch, emitter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "grab: %v\n", err)
			os.Exit(Failure)
		}
		emitter.Finalize()
		if !found {
			os.Exit(NoMatch)
		}
		os.Exit(Succeed)
	}

	for _, path := range files {
		if fileutil.IsDir(path) {
			if !opts.recursive {
				if !opts.noErrMessage {
					fmt.Fprintf(os.Stderr, "grab: %s is a directory\n", path)
				}
				continue
			}
			walkDir(path, mtch, emitter, opts)
			continue
		}

		f, err := os.Open(path)
		if err != nil {
			if !opts.noErrMessage {
				fmt.Fprintf(os.Stderr, "grab: %v\n", err)
			}
			continue
		}

		if !fileutil.IsLikelyText(f) {
			f.Close()
			continue
		}
		_, _ = f.Seek(0, io.SeekStart)

		found, err := processReader(f, path, mtch, emitter)
		f.Close()

		if err != nil {
			if !opts.noErrMessage {
				fmt.Fprintf(os.Stderr, "grab: %v\n", err)
			}
			continue
		}

		if found {
			anyFound = true
		}
	}

	emitter.Finalize()

	if !anyFound {
		os.Exit(NoMatch)
	}
}
