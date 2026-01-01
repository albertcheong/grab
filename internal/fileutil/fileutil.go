package fileutil

import (
	"io"
	"os"
)

// IsDir reports whether path refers to a directory.
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// IsELF reports whether r starts with a valid ELF magic header.
func IsELF(r io.ReaderAt) bool {
	var ident [4]byte
	n, err := r.ReadAt(ident[:], 0)
	if err != nil && err != io.EOF {
		return false
	}

	// If there's less than 4 bytes inside, than it is definitely not a file
	if n < 4 {
		return false
	}

	return ident[0] == 0x7F &&
		ident[1] == 0x45 &&
		ident[2] == 0x4c &&
		ident[3] == 0x46
}

// IsLikelyText reports whether the file appears to be text on the absence of NULL bytes in the first 512 bytes
//
// This is a heuristic and may return false positives (i.e. UTF-16 may contain `0x00`)
func IsLikelyText(r io.ReaderAt) bool {
	// Probably big enough
	var buffer [512]byte
	n, err := r.ReadAt(buffer[:], 0)
	if err != nil && err != io.EOF {
		return false
	}

	for i := range n {
		if buffer[i] == 0x00 {
			return false
		}
	}
	return true
}
