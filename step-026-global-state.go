package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

// editorConfig stores the editor's configuration, including terminal size.
var editorConfig struct {
	termWidth  int
	termHeight int
}

// initEditor initializes the editor configuration.
func initEditor() {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		die(fmt.Errorf("getting terminal size: %w", err))
	}
	editorConfig.termWidth = width
	editorConfig.termHeight = height
}

// editorDrawRows draws rows of tildes based on the terminal height.
func editorDrawRows() {
	for y := 0; y < editorConfig.termHeight; y++ {
		_, err := os.Stdout.Write([]byte("~\r\n"))
		if err != nil {
			die(fmt.Errorf("writing to stdout: %w", err))
		}
	}
}

// editorRefreshScreen clears the terminal screen and draws tildes.
func editorRefreshScreen() {
	// ANSI escape sequence to clear the entire screen.
	// "\x1b" is the escape character.
	// "[2J" clears the entire screen.
	_, err := os.Stdout.Write([]byte("\x1b[2J"))
	if err != nil {
		die(fmt.Errorf("writing to stdout: %w", err))
	}

	// ANSI escape sequence to reposition the cursor to the top-left corner.
	// "\x1b" is the escape character.
	// "[H" repositions the cursor to row 1, column 1.
	_, err = os.Stdout.Write([]byte("\x1b[H"))
	if err != nil {
		die(fmt.Errorf("writing to stdout: %w", err))
	}

	editorDrawRows()

	// ANSI escape sequence to reposition the cursor to the top-left corner.
	// "\x1b" is the escape character.
	// "[H" repositions the cursor to row 1, column 1.
	_, err = os.Stdout.Write([]byte("\x1b[H"))
	if err != nil {
		die(fmt.Errorf("writing to stdout: %w", err))
	}
}

// editorProcessKeypress waits for a keypress and handles it.
// It returns false if the editor should exit.
func editorProcessKeypress() bool {
	char, err := editorReadKey()
	if err != nil {
		if err == io.EOF {
			return false // Exit on EOF (e.g., Ctrl-D).
		}
		die(fmt.Errorf("read error: %w", err))
	}

	switch char {
	case 'q' & 0x1f: // Ctrl-Q
		// ANSI escape sequence to clear the entire screen.
		// "\x1b" is the escape character.
		// "[2J" clears the entire screen.
		_, err2 := os.Stdout.Write([]byte("\x1b[2J"))
		if err2 != nil {
			die(fmt.Errorf("writing to stdout: %w", err2))
		}

		// ANSI escape sequence to reposition the cursor to the top-left corner.
		// "\x1b" is the escape character.
		// "[H" repositions the cursor to row 1, column 1.
		_, err2 = os.Stdout.Write([]byte("\x1b[H"))
		if err2 != nil {
			die(fmt.Errorf("writing to stdout: %w", err2))
		}
		return false // Signal to exit.
	}

	return true // Signal to continue.
}

func main() {
	// Get the file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Make a copy of the original terminal state
	oldState, err := term.GetState(fd)
	if err != nil {
		die(fmt.Errorf("getting terminal state: %w", err))
	}
	// Restore the terminal state on exit
	defer term.Restore(fd, oldState)

	// Put the terminal into raw mode
	_, err = term.MakeRaw(fd)
	if err != nil {
		die(fmt.Errorf("putting terminal in raw mode: %w", err))
	}

	initEditor()

	for {
		editorRefreshScreen() // Clear the screen at the beginning of each loop iteration
		if !editorProcessKeypress() {
			break
		}
	}
}

// die panics with the given error.
// A panic will unwind the stack, run deferred functions, and exit with a stack trace.
func die(err error) {
	// ANSI escape sequence to clear the entire screen.
	// "\x1b" is the escape character.
	// "[2J" clears the entire screen.
	_, err2 := os.Stdout.Write([]byte("\x1b[2J"))
	if err2 != nil {
		// If we can't clear the screen, just print the original error
		fmt.Println("Error writing to stdout:", err2)
	}

	// ANSI escape sequence to reposition the cursor to the top-left corner.
	// "\x1b" is the escape character.
	// "[H" repositions the cursor to row 1, column 1.
	_, err2 = os.Stdout.Write([]byte("\x1b[H"))
	if err2 != nil {
		fmt.Println("Error writing to stdout:", err2)
	}

	panic(err)
}

// editorReadKey reads a single keypress from stdin and returns it.
func editorReadKey() (byte, error) {
	var b []byte = make([]byte, 1)
	n, err := os.Stdin.Read(b)
	if err != nil {
		return 0, err // Propagate error, including EOF.
	}
	if n != 1 {
		// This should not happen with a blocking read, but handle it defensively.
		return 0, io.EOF
	}
	return b[0], nil
}
