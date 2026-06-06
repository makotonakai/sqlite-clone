package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ColumnUsernameSize = 32
	ColumnEmailSize    = 255
)

const (
	IDSize       = 4
	UsernameSize = ColumnUsernameSize
	EmailSize    = ColumnEmailSize

	IDOffset       = 0
	UsernameOffset = IDOffset + IDSize
	EmailOffset    = UsernameOffset + UsernameSize

	RowSize = IDSize + UsernameSize + EmailSize
)

const (
	PageSize     = 4096
	TableMaxPages = 100
	RowsPerPage  = PageSize / RowSize
	TableMaxRows = RowsPerPage * TableMaxPages
)

type InputBuffer struct {
	Buffer string
}

type ExecuteResult int

const (
	ExecuteSuccess ExecuteResult = iota
	ExecuteTableFull
)

type MetaCommandResult int

const (
	MetaCommandSuccess MetaCommandResult = iota
	MetaCommandUnrecognizedCommand
)

type PrepareResult int

const (
	PrepareSuccess PrepareResult = iota
	PrepareSyntaxError
	PrepareUnrecognizedStatement
)

type StatementType int

const (
	StatementInsert StatementType = iota
	StatementSelect
)

type Row struct {
	ID       uint32
	Username string
	Email    string
}

type Statement struct {
	Type        StatementType
	RowToInsert Row
}

type Table struct {
	NumRows uint32
	Pages   [TableMaxPages][]byte
}

func newTable() *Table {
	return &Table{}
}

func freeTable(table *Table) {
	// No-op in Go.
	// Garbage collector handles cleanup.
}

func newInputBuffer() *InputBuffer {
	return &InputBuffer{}
}

func closeInputBuffer(inputBuffer *InputBuffer) {
	// No-op in Go.
}

func printPrompt() {
	fmt.Print("db > ")
}

func readInput(
	inputBuffer *InputBuffer,
	reader *bufio.Reader,
) {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input")
		os.Exit(1)
	}

	inputBuffer.Buffer = strings.TrimSpace(line)
}

func printRow(row *Row) {
	fmt.Printf(
		"(%d, %s, %s)\n",
		row.ID,
		row.Username,
		row.Email,
	)
}

func serializeRow(
	source *Row,
	destination []byte,
) {
	binary.LittleEndian.PutUint32(
		destination[IDOffset:IDOffset+IDSize],
		source.ID,
	)

	copy(
		destination[UsernameOffset:UsernameOffset+UsernameSize],
		[]byte(source.Username),
	)

	copy(
		destination[EmailOffset:EmailOffset+EmailSize],
		[]byte(source.Email),
	)
}

func deserializeRow(
	source []byte,
	destination *Row,
) {
	destination.ID = binary.LittleEndian.Uint32(
		source[IDOffset : IDOffset+IDSize],
	)

	destination.Username = strings.TrimRight(
		string(source[UsernameOffset:UsernameOffset+UsernameSize]),
		"\x00",
	)

	destination.Email = strings.TrimRight(
		string(source[EmailOffset:EmailOffset+EmailSize]),
		"\x00",
	)
}

func rowSlot(
	table *Table,
	rowNum uint32,
) []byte {

	pageNum := rowNum / RowsPerPage

	if table.Pages[pageNum] == nil {
		table.Pages[pageNum] = make([]byte, PageSize)
	}

	rowOffset := rowNum % RowsPerPage
	byteOffset := rowOffset * RowSize

	return table.Pages[pageNum][byteOffset : byteOffset+RowSize]
}

func doMetaCommand(
	inputBuffer *InputBuffer,
	table *Table,
) MetaCommandResult {

	if inputBuffer.Buffer == ".exit" {
		closeInputBuffer(inputBuffer)
		freeTable(table)
		os.Exit(0)
	}

	return MetaCommandUnrecognizedCommand
}

func prepareStatement(
	inputBuffer *InputBuffer,
	statement *Statement,
) PrepareResult {

	fields := strings.Fields(inputBuffer.Buffer)

	if len(fields) == 0 {
		return PrepareUnrecognizedStatement
	}

	switch fields[0] {

	case "insert":
		if len(fields) < 4 {
			return PrepareSyntaxError
		}

		id, err := strconv.ParseUint(
			fields[1],
			10,
			32,
		)

		if err != nil {
			return PrepareSyntaxError
		}

		statement.Type = StatementInsert
		statement.RowToInsert = Row{
			ID:       uint32(id),
			Username: fields[2],
			Email:    fields[3],
		}

		return PrepareSuccess

	case "select":
		statement.Type = StatementSelect
		return PrepareSuccess
	}

	return PrepareUnrecognizedStatement
}

func executeInsert(
	statement *Statement,
	table *Table,
) ExecuteResult {

	if table.NumRows >= TableMaxRows {
		return ExecuteTableFull
	}

	serializeRow(
		&statement.RowToInsert,
		rowSlot(table, table.NumRows),
	)

	table.NumRows++

	return ExecuteSuccess
}

func executeSelect(
	statement *Statement,
	table *Table,
) ExecuteResult {

	var row Row

	for i := uint32(0); i < table.NumRows; i++ {
		deserializeRow(
			rowSlot(table, i),
			&row,
		)

		printRow(&row)
	}

	return ExecuteSuccess
}

func executeStatement(
	statement *Statement,
	table *Table,
) ExecuteResult {

	switch statement.Type {

	case StatementInsert:
		return executeInsert(statement, table)

	case StatementSelect:
		return executeSelect(statement, table)
	}

	return ExecuteSuccess
}

func main() {

	table := newTable()
	inputBuffer := newInputBuffer()
	reader := bufio.NewReader(os.Stdin)

	for {
		printPrompt()
		readInput(inputBuffer, reader)

		if len(inputBuffer.Buffer) > 0 &&
			inputBuffer.Buffer[0] == '.' {

			switch doMetaCommand(
				inputBuffer,
				table,
			) {

			case MetaCommandSuccess:
				continue

			case MetaCommandUnrecognizedCommand:
				fmt.Printf(
					"Unrecognized command '%s'\n",
					inputBuffer.Buffer,
				)
				continue
			}
		} else if len(inputBuffer.Buffer) == 0 {
				continue
		}

		var statement Statement

		switch prepareStatement(
			inputBuffer,
			&statement,
		) {

		case PrepareSuccess:

		case PrepareSyntaxError:
			fmt.Println(
				"Syntax error. Could not parse statement.",
			)
			continue

		case PrepareUnrecognizedStatement:
			fmt.Printf(
				"Unrecognized keyword at start of '%s'.\n",
				inputBuffer.Buffer,
			)
			continue
		}

		switch executeStatement(
			&statement,
			table,
		) {

		case ExecuteSuccess:
			fmt.Println("Executed.")

		case ExecuteTableFull:
			fmt.Println("Error: Table full.")
		}
	}
}
