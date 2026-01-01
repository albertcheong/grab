package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "testdata"

func TestIsDir(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "directory exists",
			path:     filepath.Join(testDataDir, "dummyDir"),
			expected: true,
		},
		{
			name:     "file is not directory",
			path:     filepath.Join(testDataDir, "dummy"),
			expected: false,
		},
		{
			name:     "another file is not directory",
			path:     filepath.Join(testDataDir, "utf8.txt"),
			expected: false,
		},
		{
			name:     "path does not exist",
			path:     filepath.Join(testDataDir, "nonexistent"),
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDir(tt.path)
			if result != tt.expected {
				t.Errorf("IsDir(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsELF(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "ELF binary",
			path:     filepath.Join(testDataDir, "dummy"),
			expected: true,
		},
		{
			name:     "Windows executable (not ELF)",
			path:     filepath.Join(testDataDir, "dummy.exe"),
			expected: false,
		},
		{
			name:     "Python script (not ELF)",
			path:     filepath.Join(testDataDir, "dummy.py"),
			expected: false,
		},
		{
			name:     "UTF-8 text file (not ELF)",
			path:     filepath.Join(testDataDir, "utf8.txt"),
			expected: false,
		},
		{
			name:     "PNG image (not ELF)",
			path:     filepath.Join(testDataDir, "dummy.png"),
			expected: false,
		},
		{
			name:     "HTML file (not ELF)",
			path:     filepath.Join(testDataDir, "dummy.html"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.path)
			if err != nil {
				t.Fatalf("Failed to open file %q: %v", tt.path, err)
			}
			defer f.Close()

			result := IsELF(f)
			if result != tt.expected {
				t.Errorf("IsELF(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsELF_EdgeCases(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		// Create a temporary empty file
		tmpFile, err := os.CreateTemp("", "empty-*.bin")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		result := IsELF(tmpFile)
		if result != false {
			t.Errorf("IsELF(empty file) = %v, expected false", result)
		}
	})

	t.Run("file with less than 4 bytes", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "small-*.bin")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write only 3 bytes
		if _, err := tmpFile.Write([]byte{0x7F, 0x45, 0x4c}); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		f, _ := os.Open(tmpFile.Name())
		defer f.Close()

		result := IsELF(f)
		if result != false {
			t.Errorf("IsELF(3 bytes) = %v, expected false", result)
		}
	})
}

func TestIsLikelyText(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "UTF-8 text file",
			path:     filepath.Join(testDataDir, "utf8.txt"),
			expected: true,
		},
		{
			name:     "Python script",
			path:     filepath.Join(testDataDir, "dummy.py"),
			expected: true,
		},
		{
			name:     "HTML file",
			path:     filepath.Join(testDataDir, "dummy.html"),
			expected: true,
		},
		{
			name:     "UTF-16 file (contains null bytes)",
			path:     filepath.Join(testDataDir, "utf16.txt"),
			expected: false,
		},
		{
			name:     "ELF binary (contains null bytes)",
			path:     filepath.Join(testDataDir, "dummy"),
			expected: false,
		},
		{
			name:     "Windows executable (contains null bytes)",
			path:     filepath.Join(testDataDir, "dummy.exe"),
			expected: false,
		},
		{
			name:     "PNG image (binary)",
			path:     filepath.Join(testDataDir, "dummy.png"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.path)
			if err != nil {
				t.Fatalf("Failed to open file %q: %v", tt.path, err)
			}
			defer f.Close()

			result := IsLikelyText(f)
			if result != tt.expected {
				t.Errorf("IsLikelyText(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsLikelyText_EdgeCases(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "empty-*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		result := IsLikelyText(tmpFile)
		if result != true {
			t.Errorf("IsLikelyText(empty file) = %v, expected true", result)
		}
	})

	t.Run("file with null byte at position 0", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "null-*.bin")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.Write([]byte{0x00, 'h', 'e', 'l', 'l', 'o'}); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		f, _ := os.Open(tmpFile.Name())
		defer f.Close()

		result := IsLikelyText(f)
		if result != false {
			t.Errorf("IsLikelyText(null at start) = %v, expected false", result)
		}
	})

	t.Run("file with null byte at position 511", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "null-end-*.bin")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write 511 bytes of 'A' followed by a null byte
		data := make([]byte, 512)
		for i := range 511 {
			data[i] = 'A'
		}
		data[511] = 0x00

		if _, err := tmpFile.Write(data); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		f, _ := os.Open(tmpFile.Name())
		defer f.Close()

		result := IsLikelyText(f)
		if result != false {
			t.Errorf("IsLikelyText(null at byte 511) = %v, expected false", result)
		}
	})

	t.Run("file larger than 512 bytes, all text", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "large-*.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write 1000 bytes of text
		data := make([]byte, 1000)
		for i := range data {
			data[i] = 'A'
		}

		if _, err := tmpFile.Write(data); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		f, _ := os.Open(tmpFile.Name())
		defer f.Close()

		result := IsLikelyText(f)
		if result != true {
			t.Errorf("IsLikelyText(large text file) = %v, expected true", result)
		}
	})
}
