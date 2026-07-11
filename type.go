package main

import (
	"os"
	"fmt"
	"strings"
    "strconv"
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
    PREPARE_SYNTAX_ERROR
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

type Row struct {
    ID uint32
    UserName string
    Email string
}

type Statement struct {
    Type StatementType
    RowToInsert Row
}

func PrepareStatement(line string, statement *Statement) PrepareResult {

    fields := strings.Fields(line)

    if len(fields) == 0 {
        return PREPARE_UNRECOGNIZED_COMMAND
    }

    switch fields[0] {
    case "insert":
        fields := strings.Fields(line)
        if len(fields) < 4 {
            return PREPARE_SYNTAX_ERROR
        }

        id, err := strconv.Atoi(fields[1])
        if err != nil {
            return PREPARE_SYNTAX_ERROR
        }

        statement.Type = STATEMENT_INSERT
        statement.RowToInsert = Row{
            ID: uint32(id),
            UserName: fields[2],
            Email: fields[3],
        }

        return PREPARE_SUCCESS

    case "select":
        statement.Type = STATEMENT_SELECT
        return PREPARE_SUCCESS

    default:
        return PREPARE_UNRECOGNIZED_COMMAND
    }

}

func ExecuteInsert(statement *Statement, table *Table) ExecuteResult {

    if len(table.Rows) >= TABLE_MAX_ROWS {
        return EXECUTE_TABLE_FULL
    }

    table.Rows = append(table.Rows, statement.RowToInsert)

    return EXECUTE_SUCCESS

}

func ExecuteSelect(table *Table) ExecuteResult {

    for _, row := range table.Rows {
        fmt.Printf("(%d, %s, %s)\n",
            row.ID,
            row.UserName,
            row.Email,
        )
    }

    return EXECUTE_SUCCESS
}

func ExecuteStatement(statement *Statement, table *Table) ExecuteResult {

    switch statement.Type {

    case STATEMENT_INSERT:
        return ExecuteInsert(statement, table)
        
    case STATEMENT_SELECT:
        return ExecuteSelect(table)
    }

    return EXECUTE_SUCCESS
    
}

type ExecuteResult int

const (
    EXECUTE_SUCCESS = iota
    EXECUTE_TABLE_FULL
)

type Table struct {
    Rows []Row
}

func NewTable() *Table {
    return &Table{
        Rows: make([]Row, 0),
    }
}
