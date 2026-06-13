package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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
	PageSize      = 4096
	TableMaxPages = 100
	RowsPerPage   = PageSize / RowSize
	TableMaxRows  = RowsPerPage * TableMaxPages
)

type InputBuffer struct {
	Buffer string
}

type ExecuteResult int

const (
	ExecuteSuccess ExecuteResult = iota
	ExecuteTableFull
	ExecuteDuplicateKey
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

type Pager struct {
	File       *os.File
	FileLength uint32
	Pages      [TableMaxPages][]byte
}

type Table struct {
	NumRows uint32
	Pager   *Pager
}

func printPrompt() {
	fmt.Print("db > ")
}

func readInput(inputBuffer *InputBuffer, reader *bufio.Reader) {
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Println("Error reading input")
		os.Exit(1)
	}

	inputBuffer.Buffer = strings.TrimSpace(line)
}

func newInputBuffer() *InputBuffer {
	return &InputBuffer{}
}

func serializeRow(source *Row, destination []byte) {
	binary.LittleEndian.PutUint32(destination[IDOffset:IDOffset+IDSize], source.ID)
	copy(destination[UsernameOffset:UsernameOffset+UsernameSize], []byte(source.Username))
	copy(destination[EmailOffset:EmailOffset+EmailSize], []byte(source.Email))
}

func deserializeRow(source []byte, destination *Row) {
	destination.ID = binary.LittleEndian.Uint32(source[IDOffset : IDOffset+IDSize])
	destination.Username = strings.TrimRight(
		string(source[UsernameOffset:UsernameOffset+UsernameSize]), "\x00")
	destination.Email = strings.TrimRight(
		string(source[EmailOffset:EmailOffset+EmailSize]), "\x00")
}

func printRow(row *Row) {
	fmt.Printf("(%d, %s, %s)\n", row.ID, row.Username, row.Email)
}

func pagerOpen(filename string) *Pager {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Unable to open file")
		os.Exit(1)
	}

	info, err := file.Stat()
	if err != nil {
		fmt.Println("Unable to stat file")
		os.Exit(1)
	}

	return &Pager{
		File:       file,
		FileLength: uint32(info.Size()),
	}
}

func dbOpen(filename string) *Table {
	pager := pagerOpen(filename)

	return &Table{
		Pager:   pager,
		NumRows: pager.FileLength / RowSize,
	}
}

func getPage(pager *Pager, pageNum uint32) []byte {
	if pageNum >= TableMaxPages {
		fmt.Printf("Tried to fetch page number out of bounds. %d >= %d\n",
			pageNum, TableMaxPages)
		os.Exit(1)
	}

	if pager.Pages[pageNum] == nil {
		page := make([]byte, PageSize)

		numPages := pager.FileLength / PageSize
		if pager.FileLength%PageSize != 0 {
			numPages++
		}

		if pageNum < numPages {
			_, err := pager.File.ReadAt(page, int64(pageNum*PageSize))
			if err != nil && err != io.EOF {
				fmt.Println("Error reading file:", err)
				os.Exit(1)
			}
		}

		pager.Pages[pageNum] = page
	}

	return pager.Pages[pageNum]
}

func rowSlot(table *Table, rowNum uint32) []byte {
	pageNum := rowNum / RowsPerPage
	page := getPage(table.Pager, pageNum)

	rowOffset := rowNum % RowsPerPage
	byteOffset := rowOffset * RowSize

	return page[byteOffset : byteOffset+RowSize]
}

func pagerFlush(pager *Pager, pageNum uint32, size uint32) {
	page := pager.Pages[pageNum]
	if page == nil {
		fmt.Println("Tried to flush null page")
		os.Exit(1)
	}

	_, err := pager.File.WriteAt(page[:size], int64(pageNum*PageSize))
	if err != nil {
		fmt.Println("Error writing:", err)
		os.Exit(1)
	}
}

func dbClose(table *Table) {
	pager := table.Pager

	numFullPages := table.NumRows / RowsPerPage

	for i := uint32(0); i < numFullPages; i++ {
		if pager.Pages[i] == nil {
			continue
		}
		pagerFlush(pager, i, PageSize)
		pager.Pages[i] = nil
	}

	numAdditionalRows := table.NumRows % RowsPerPage
	if numAdditionalRows > 0 {
		pageNum := numFullPages
		if pager.Pages[pageNum] != nil {
			pagerFlush(pager, pageNum, numAdditionalRows*RowSize)
			pager.Pages[pageNum] = nil
		}
	}

	if err := pager.File.Close(); err != nil {
		fmt.Println("Error closing db file")
		os.Exit(1)
	}
}

func doMetaCommand(inputBuffer *InputBuffer, table *Table) MetaCommandResult {
	if inputBuffer.Buffer == ".exit" {
		dbClose(table)
		os.Exit(0)
	}
	return MetaCommandUnrecognizedCommand
}

func prepareStatement(inputBuffer *InputBuffer, statement *Statement) PrepareResult {
	fields := strings.Fields(inputBuffer.Buffer)

	if len(fields) == 0 {
		return PrepareUnrecognizedStatement
	}

	switch fields[0] {
	case "insert":
		if len(fields) < 4 {
			return PrepareSyntaxError
		}

		id, err := strconv.ParseUint(fields[1], 10, 32)
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

func executeInsert(statement *Statement, table *Table) ExecuteResult {
	if table.NumRows >= TableMaxRows {
		return ExecuteTableFull
	}

	idToInsert := statement.RowToInsert.ID

	var row Row
	for i := uint32(0); i < table.NumRows; i++ {
		deserializeRow(rowSlot(table, i), &row)
		if row.ID == idToInsert {
			return ExecuteDuplicateKey
		}
	}

	serializeRow(&statement.RowToInsert, rowSlot(table, table.NumRows))
	table.NumRows++

	return ExecuteSuccess
}

func executeSelect(table *Table) ExecuteResult {
	var row Row

	for i := uint32(0); i < table.NumRows; i++ {
		deserializeRow(rowSlot(table, i), &row)
		printRow(&row)
	}

	return ExecuteSuccess
}

func executeStatement(statement *Statement, table *Table) ExecuteResult {
	switch statement.Type {
	case StatementInsert:
		return executeInsert(statement, table)
	case StatementSelect:
		return executeSelect(table)
	}
	return ExecuteSuccess
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Must supply a database filename.")
		os.Exit(1)
	}

	table := dbOpen(os.Args[1])
	inputBuffer := newInputBuffer()
	reader := bufio.NewReader(os.Stdin)

	for {
		printPrompt()
		readInput(inputBuffer, reader)

		if len(inputBuffer.Buffer) == 0 {
			continue
		}

		if strings.HasPrefix(inputBuffer.Buffer, ".") {
			switch doMetaCommand(inputBuffer, table) {
			case MetaCommandUnrecognizedCommand:
				fmt.Printf("Unrecognized command '%s'\n", inputBuffer.Buffer)
				continue
			}
		}

		var statement Statement

		switch prepareStatement(inputBuffer, &statement) {
		case PrepareSyntaxError:
			fmt.Println("Syntax error. Could not parse statement.")
			continue
		case PrepareUnrecognizedStatement:
			fmt.Printf("Unrecognized keyword at start of '%s'.\n",
				inputBuffer.Buffer)
			continue
		}

		switch executeStatement(&statement, table) {
		case ExecuteSuccess:
			fmt.Println("Executed.")
		case ExecuteTableFull:
			fmt.Println("Error: Table full.")
		case ExecuteDuplicateKey:
			fmt.Println("Error: Duplicate key.")
		}
	}
}
