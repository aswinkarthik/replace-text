package fs

import (
	"fmt"
	"io"
	"os"
)

// Fs is an abstraction to access the filesystem
type Fs interface {
	// IsFile returns true if the path is a file
	// Returns false if its a directory, or does not exists, or on IO Error
	IsFile(path string) bool

	// IsDir returns true if the path is a directory.
	// Returns false if its a file, or does not exists, or on IO Error
	IsDir(path string) bool

	// DevNull returns a noop writer with no side effects
	// An OSFs should return /dev/null
	DevNull() io.Writer

	// Open will open a File handle to the path. Returns error
	// if file is missing. The file is opened as a readonly file.
	Open(path string) (ReadOnlyFile, error)

	// Create will open a file for writing with the specified file mode.
	// Contents will be truncated on open.
	// It will error out if file already exists.
	Create(path string, mode os.FileMode) (WritableFile, error)

	// Exists returns true if is a valid file or directory.
	Exists(path string) (bool, error)

	// FileMode will return the file mode of the given path.
	FileMode(path string) (os.FileMode, error)
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

func (f *osFs) Open(path string) (ReadOnlyFile, error) {
	return os.Open(path)
}

func (f *osFs) Create(path string, mode os.FileMode) (WritableFile, error) {
	exists, err := f.Exists(path)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}

	if exists {
		return nil, fmt.Errorf("cannot create file %s as it already exists", path)
	}

	return os.OpenFile(path, os.O_CREATE, mode)
}

func (f *osFs) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if err == os.ErrNotExist {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (f *osFs) FileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}
