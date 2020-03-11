package replacer

import (
	"fmt"
	"io"
)

// Replacer is the struct responsible for doing IO operations
type Replacer struct {
	stateMachines *StateMachines
}

// NewReplacer is a constructor for creating Replacer struct.
// Accepts replacements and initializes state machines with trie
func NewReplacer(replacements map[string]string) (*Replacer, error) {
	root := NewNode()
	for k, v := range replacements {
		if err := root.Put(k, v); err != nil {
			return nil, fmt.Errorf("error creating replacer: %v", err)
		}
	}

	sm := NewStateMachines(root)

	return &Replacer{sm}, nil
}

// Run accepts a reader and writer.
// Data from reader is copied into writer.
// While doing so, it replaces all found matches with replace value.
//
// It replaces all matches in 2 passes.
//
// The first pass is to move all the state machines.
// Second pass to make use of all the terminal nodes to make the replacements
// in the writer.
func (r *Replacer) Run(reader io.ReadSeeker, writer io.Writer) error {
	const bufferSize = 8000

	return r.run(bufferSize, reader, writer)
}

func (r *Replacer) run(bufferSize int, reader io.ReadSeeker, writer io.Writer) error {
	// Construct the state machines first
	for position, readBuffer := int64(0), make([]byte, bufferSize); true; {
		n, err := reader.Read(readBuffer)
		if n > 0 {
			for _, b := range readBuffer[:n] {
				r.stateMachines.Accept(b, position)
				position++
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error finding matches: %v", err)
		}
	}

	if len(r.stateMachines.TerminalMachines) == 0 {
		return fmt.Errorf("no matches found")
	}

	// Reset to beginning of file
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking to start of file: %v", err)
	}

	// n represents total bytes read from reader
	var n int64
	for _, m := range r.stateMachines.TerminalMachines {
		// Copy till first match
		if _, err := io.CopyN(writer, reader, m.StartPosition-n); err != nil {
			return fmt.Errorf("error copying data from source to destination: %v", err)
		}

		// Print the replacement string
		if _, err := writer.Write([]byte(m.ReplaceWith)); err != nil {
			return fmt.Errorf("error writing replaced strings: %v", err)
		}

		// Seek to the end position and move on to next match
		if _, err := reader.Seek(m.EndPosition+1, io.SeekStart); err != nil {
			return fmt.Errorf("error seeking to next location: %v", err)
		}

		// Update total bytes read
		n = m.EndPosition + 1
	}

	// Copy remaining data.
	_, err := io.Copy(writer, reader)
	return err
}
