package main

// TODO: Let’s make a function for low-level keypress reading, and another function for mapping keypresses to editor operations. We’ll also stop printing out keypresses at this point.
// char editorReadKey() {
//   int nread;
//   char c;
//   while ((nread = read(STDIN_FILENO, &c, 1)) != 1) {
//     if (nread == -1 && errno != EAGAIN) die("read");
//   }
//   return c;
// }
// /*** input ***/
// void editorProcessKeypress() {
//   char c = editorReadKey();
//   switch (c) {
//     case CTRL_KEY('q'):
//       exit(0);
//       break;
//   }
// }
// /*** init ***/
// int main() {
//   enableRawMode();
//   while (1) {
//     editorProcessKeypress();
//   }
//   return 0;
// }

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

// die panics with the given error.
// A panic will unwind the stack, run deferred functions, and exit with a stack trace.
func die(err error) {
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

	for {
		if !editorProcessKeypress() {
			break
		}
	}
}
