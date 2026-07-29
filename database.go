package main

import (
	"strings"
  "encoding/binary"
)

type Table struct {
    Pager *Pager
    NumRows uint32
}

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



