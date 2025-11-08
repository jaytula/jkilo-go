package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		c, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			// Handle other errors if necessary, or just break
			break
		}
		if c == 'q' {
			break // 'q' character encountered
		}
	}
	os.Exit(0)
}
