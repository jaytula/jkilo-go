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
