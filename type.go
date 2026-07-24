package main

import (
    "io"
	"os"
	"fmt"
	"strings"
    "strconv"
    "syscall"
    "encoding/binary"
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


func DoMetaCommand(line string, table *Table) MetaCommandResult {

    if line == ".exit" {
        DBClose(table)
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

    if table.NumRows >= TABLE_MAX_ROWS {
        return EXECUTE_TABLE_FULL
    }

    cursor := EndTable(table)
    SerializeRow(&statement.RowToInsert, CursorValue(cursor))

    table.NumRows++

    return EXECUTE_SUCCESS
}

func ExecuteSelect(table *Table) ExecuteResult {

    cursor := StartTable(table)

    var row Row

    for !cursor.EndOfTable {

        DeserializeRow(CursorValue(cursor), &row)

        fmt.Printf("(%d, %s, %s)\n",
            row.ID,
            row.UserName,
            row.Email,
        )

        CursorAdvance(cursor)
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
    Pager *Pager
    NumRows uint32
}

// func NewTable() *Table {
//     return &Table{
//         P
//     }
// }

func DBOpen(fileName string) *Table {

    p := PagerOpen(fileName)

    return &Table {
        Pager: p,
        NumRows: uint32(p.FileLength / ROW_SIZE),
    }

}

func DBClose(table *Table) {

    pager := table.Pager

    numFullPages := table.NumRows / ROWS_PER_PAGE

    for i := uint32(0); i < numFullPages; i++ {

        if pager.Pages[i] == nil {
            continue
        }

        PagerFlush(pager, i, PAGE_SIZE)

        pager.Pages[i] = nil
    }

    numAdditionalRows := table.NumRows % ROWS_PER_PAGE

    if numAdditionalRows > 0 {

        pageNum := numFullPages

        if pager.Pages[pageNum] != nil {

            PagerFlush(
                pager,
                pageNum,
                numAdditionalRows*ROW_SIZE,
            )

            pager.Pages[pageNum] = nil
        }
    }

    pager.File.Close()
}

func PagerOpen(fileName string) *Pager {
    
    file, err := os.OpenFile(
        fileName, 
        os.O_RDWR|os.O_CREATE|syscall.S_IWUSR|syscall.S_IRUSR, 
        0644,
    )

    if err != nil {
        fmt.Println(err)
    }

    stat, err := file.Stat()
    if err != nil {
        fmt.Println(err)
    }

    p := Pager{
        File: file,
        FileLength: stat.Size(),
    }

    return &p
}

func PagerFlush(pager *Pager, pageNum uint32, size uint32) {

    if pager.Pages[pageNum] == nil {
        fmt.Printf("tried to flush nil page\n")
    }

    _, err := pager.File.Seek(
        int64(pageNum)*PAGE_SIZE,
        0,
    )

    if err != nil {
        fmt.Println(err)
    }

    _, err = pager.File.Write(
        pager.Pages[pageNum][:size],
    )

    if err != nil {
        fmt.Println(err)
    }
}

func GetPage(pager *Pager, pageNum uint32) []byte {

    if pageNum >= TABLE_MAX_PAGES {
        fmt.Printf("Tried to fetch page number out of bounds. %d > %d\n", 
            pageNum,
            TABLE_MAX_PAGES,
        );
        os.Exit(1)
    }

    if pager.Pages[pageNum] == nil {
        page := make([]byte, PAGE_SIZE)
        numPages := pager.FileLength / PAGE_SIZE

        if pager.FileLength % PAGE_SIZE != 0 {
            numPages++
        }

        if int64(pageNum) < numPages {
            pager.File.Seek(int64(pageNum)*PAGE_SIZE, 0)

            _, err := io.ReadFull(pager.File, page)
            if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
                fmt.Println(err)
                os.Exit(1)
            }
        }

        pager.Pages[pageNum] = page
    }

    return pager.Pages[pageNum]
}

func CursorValue(cursor *Cursor) []byte {

    numRows := cursor.NumRows
    pageNum := numRows / ROWS_PER_PAGE

    p := GetPage(cursor.Table.Pager, pageNum)

    rowOffset := numRows % ROWS_PER_PAGE
    byteOffset := rowOffset * ROW_SIZE

    return p[byteOffset: byteOffset+ROW_SIZE]
}

func CursorAdvance(cursor *Cursor) {

    cursor.NumRows = cursor.NumRows + 1

    if cursor.NumRows >= cursor.Table.NumRows {
        cursor.EndOfTable = true
    }

}

type Pager struct {
    File *os.File
    FileLength int64
    Pages [TABLE_MAX_PAGES][]byte
}

func SerializeRow(row *Row, destination []byte) {

    // id
    binary.LittleEndian.PutUint32(
        destination[ID_OFFSET:ID_OFFSET+ID_SIZE],
        row.ID,
    )

    // username
    copy(
        destination[USERNAME_OFFSET:USERNAME_OFFSET+USERNAME_SIZE],
        []byte(row.UserName),
    )

    // email
    copy(
        destination[EMAIL_OFFSET:EMAIL_OFFSET+EMAIL_SIZE],
        []byte(row.Email),
    )
}

func DeserializeRow(source []byte, row *Row) {

    row.ID = binary.LittleEndian.Uint32(
        source[ID_OFFSET:ID_OFFSET+ID_SIZE],
    )

    username := source[USERNAME_OFFSET : USERNAME_OFFSET+USERNAME_SIZE]
    email := source[EMAIL_OFFSET : EMAIL_OFFSET+EMAIL_SIZE]

    row.UserName = string(username)
    row.Email = string(email)

    // Remove trailing '\0'
    row.UserName = strings.TrimRight(row.UserName, "\x00")
    row.Email = strings.TrimRight(row.Email, "\x00")
}

type Cursor struct {
    Table *Table
    NumRows uint32
    EndOfTable bool
}

func StartTable(table *Table) *Cursor {

    var hasZeroRows bool
    if table.NumRows == 0 {
        hasZeroRows = true
    } else {
        hasZeroRows = false
    }

    return &Cursor{
        Table: table,
        NumRows: 0,
        EndOfTable: hasZeroRows,
    }

}

func EndTable(table *Table) *Cursor {

    return &Cursor{
        Table: table,
        NumRows: table.NumRows,
        EndOfTable: true,
    }

}



