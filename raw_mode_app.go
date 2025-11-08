package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	// Get the file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Make a copy of the original terminal state
	oldState, err := term.GetState(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal state: %v\n", err)
		os.Exit(1)
	}
	// Restore the terminal state on exit
	defer term.Restore(fd, oldState)

	// Put the terminal into raw mode
	_, err = term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error putting terminal in raw mode: %v\n", err)
		os.Exit(1)
	}

	// Read from stdin in a loop
	var b []byte = make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			break
		}
		if b[0] == 'q' {
			break
		}
	}
}
