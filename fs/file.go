package fs

import "io"

// ReadOnlyFile represents a file in Fs
// on which only read operations can be performed.
type ReadOnlyFile interface {
	io.Reader
	io.Closer
	io.Seeker
}

// WritableFile represents a file in a filesystem
// on which only write operations can be performed.
type WritableFile interface {
	io.Writer
	io.Closer
}
