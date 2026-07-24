package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

const (
    COLUMN_USERNAME_SIZE = 32
    COLUMN_EMAIL_SIZE    = 255

    ID_SIZE       = 4 // uint32
    USERNAME_SIZE = COLUMN_USERNAME_SIZE
    EMAIL_SIZE    = COLUMN_EMAIL_SIZE

    ID_OFFSET       = 0
    USERNAME_OFFSET = ID_OFFSET + ID_SIZE
    EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE

    ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

// TODO: Estimate the cache (memory) size / Why does the file take 4k?
// Data that are not on the cache will be written to the disk, which will take 1 million times.

const TABLE_MAX_PAGES=100
const PAGE_SIZE=4096
const ROWS_PER_PAGE=PAGE_SIZE/ROW_SIZE
const TABLE_MAX_ROWS=ROWS_PER_PAGE * TABLE_MAX_PAGES


func PrintPrompt() {
    fmt.Print("db > ")
}

func main() {

    r := bufio.NewReaderSize(os.Stdin, 1024*1024)
    // table := NewTable()

    if len(os.Args) < 2 {
        fmt.Printf("Must supply a database filename.\n")
        os.Exit(1)
    }

    fn := os.Args[1]
    t := DBOpen(fn)

    for {

        PrintPrompt()

        l, err := r.ReadString('\n')
        if err != nil {
            return
        }

        l = strings.TrimSpace(l)

        if strings.HasPrefix(l, ".") {
            switch DoMetaCommand(l, t) {
            case META_COMMAND_SUCCESS:
                continue
            case META_COMMAND_UNRECOGNIZED_COMMAND:
                fmt.Printf("Unrecognized command '%s'\n", l)
                continue
            }
        }

        var s Statement
        switch PrepareStatement(l, &s) {
        case PREPARE_SUCCESS:
            break
        case PREPARE_UNRECOGNIZED_COMMAND:
            fmt.Printf("Unrecognized keyword at start of %s.\n", l)
            continue
        }

        ExecuteStatement(&s, t)
        fmt.Printf("Executed.\n")
        
    }
		
}
