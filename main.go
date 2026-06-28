package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func main() {

    reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

    for {

        fmt.Print("db > ")

        line, err := reader.ReadString('\n')
        if err != nil {
            return
        }

        line = strings.TrimSpace(line)

        if line == "exit" {
            os.Exit(0)
        } else {
            fmt.Printf("Unrecognized command: %s\n", line)
        }
    }
		
}
