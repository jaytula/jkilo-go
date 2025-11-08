package main

// TODO: Last chapter we saw that the Ctrl key combined with the alphabetic keys seemed to map to bytes 1–26. We can use this to detect Ctrl key combinations and map them to different operations in our editor. We’ll start by mapping Ctrl-Q to the quit operation.


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
		n, err := os.Stdin.Read(b)
		if err != nil {
			break
		}
		if n > 0 {
			char := b[0]
			if char >= 32 && char <= 126 { // Printable ASCII characters
				fmt.Printf("%d ('%c')\r\n", char, char)
			} else {
				fmt.Printf("%d\r\n", char)
			}
		}
		if b[0] == 'q'&0x1f {
			break
		}
	}
}
