package main

import (
	"os"
	"fmt"
	"strings"
)

type MetaCommandResult int

const (
    META_COMMAND_SUCCESS = iota
    META_COMMAND_UNRECOGNIZED_COMMAND
)

type PrepareResult int

const (
    PREPARE_SUCCESS = iota
    PREPARE_UNRECOGNIZED_COMMAND
)


func DoMetaCommand(line string) MetaCommandResult {

    if line == ".exit" {
        os.Exit(0)
    } else {
        return META_COMMAND_UNRECOGNIZED_COMMAND
    }

    return META_COMMAND_SUCCESS
}

type StatementType int

const (
    STATEMENT_INSERT = iota
    STATEMENT_SELECT
)

type Statement struct {
    Type StatementType
}

func PrepareStatement(line string, statement Statement) PrepareResult {

    command := strings.Split(line, " ")

    if command[0] == "insert" {
        statement.Type = STATEMENT_INSERT
        return PREPARE_SUCCESS

    } else if command[0] == "select" {
        statement.Type = STATEMENT_SELECT
        return PREPARE_SUCCESS

    } else {
        return PREPARE_UNRECOGNIZED_COMMAND
    }

}

func ExecuteStatement(statement Statement) {

    switch statement.Type {

    case STATEMENT_INSERT:
        fmt.Printf("This is where we would do an insert.\n");
        break;

    case STATEMENT_SELECT:
        fmt.Printf("This is where we would do a select.\n");
        break;

    }
    
}
