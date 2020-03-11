package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		errStr := err.Error()
		lines := strings.Split(errStr, "\n")
		for _, line := range lines {
			_, _ = fmt.Fprintln(os.Stderr, "replace-text:", line)
		}
		os.Exit(1)
	}
}

func run() error {
	return fmt.Errorf("not implemented")
}
