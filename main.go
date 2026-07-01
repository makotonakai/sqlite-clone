package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func PrintPrompt() {
    fmt.Print("db > ")
}

func main() {

    reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

    for {

        PrintPrompt()

        line, err := reader.ReadString('\n')
        if err != nil {
            return
        }

        line = strings.TrimSpace(line)

        if strings.HasPrefix(line, ".") {
            switch DoMetaCommand(line) {
            case META_COMMAND_SUCCESS:
                continue
            case META_COMMAND_UNRECOGNIZED_COMMAND:
                fmt.Printf("Unrecognized command '%s'\n", line)
                continue
            }
        }

        var statement Statement
        switch PrepareStatement(line, statement) {
        case PREPARE_SUCCESS:
            break
        case PREPARE_UNRECOGNIZED_COMMAND:
            fmt.Printf("Unrecognized keyword at start of %s.\n", line)
            continue
        }

        ExecuteStatement(statement)
        fmt.Printf("Executed.\n")
        
    }
		
}
