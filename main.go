package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type InputBuffer struct {
	buffer      string
	inputLength in
}

func newInputBuffer() *InputBuffer {
	return &InputBuffer{}
}

func printPrompt() {
	fmt.Print("db > ")
}

func readInput(inputBuffer *InputBuffer, reader *bufio.Reader) {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input")
		os.Exit(1)
	}

	// Ignore trailing newline
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r") // Handle Windows line endings

	inputBuffer.buffer = line
	inputBuffer.inputLength = len(line)
}

func closeInputBuffer(inputBuffer *InputBuffer) {
	// No-op in Go. Garbage collector handles memory cleanup.
}

func main() {
	inputBuffer := newInputBuffer()
	reader := bufio.NewReader(os.Stdin)

	for {
		printPrompt()
		readInput(inputBuffer, reader)

		if inputBuffer.buffer == ".exit" {
			closeInputBuffer(inputBuffer)
			os.Exit(0)
		} else {
			fmt.Printf("Unrecognized command '%s'.\n", inputBuffer.buffer)
		}
	}
}
