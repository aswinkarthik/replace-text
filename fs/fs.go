package fs

import (
	"io"
	"os"
)

// Fs is an abstraction to access the filesystem
type Fs interface {
	// IsFile returns true if the path is a file
	IsFile(path string) bool
	// IsDir returns true if the path is a directory
	IsDir(path string) bool
	// DevNull returns a noop writer with no side effects
	// An OSFs should return /dev/null
	DevNull() io.Writer
}

var _ Fs = (*osFs)(nil)

type osFs struct{}

// NewOsFs can be used to create an OS FS
// that uses the actual filesystem implementation
func NewOsFs() Fs {
	return &osFs{}
}

func (f *osFs) IsFile(path string) bool {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileinfo.IsDir() {
		return false
	}

	return true
}

func (f *osFs) IsDir(path string) bool {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileinfo.IsDir() {
		return true
	}

	return false
}

func (f *osFs) DevNull() io.Writer {
	devNull, _ := os.Open(os.DevNull)
	return devNull
}
