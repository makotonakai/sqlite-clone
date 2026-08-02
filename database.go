package main

import (
    "strings"
    "encoding/binary"
)

type Table struct {
    Pager *Pager
    NumRows uint32
    RootPageNum uint32
}

func DBOpen(fileName string) *Table {

    p := PagerOpen(fileName)

    t := &Table {
        Pager: p,
        RootPageNum: 0,
    }

    if p.NumPages == 0 {
        // Root node
        rn := GetPage(p, 0)
        initializeLeafNode(rn)
    }

    return t

}

func DBClose(table *Table) {

    pager := table.Pager

    // numFullPages := table.NumRows / ROWS_PER_PAGE

    for i := 0; i < int(pager.NumPages); i++ {

        if pager.Pages[i] == nil {
            continue
        }

        PagerFlush(pager, uint32(i))

        pager.Pages[i] = nil
    }

    pager.File.Close()
}



func CursorValue(cursor *Cursor) []byte {

    pageNum := cursor.PageNum

    p := GetPage(cursor.Table.Pager, pageNum)

    return leafNodeValue(p, cursor.CellNum)
}

func CursorAdvance(cursor *Cursor) {

    // Page number
    pn := cursor.PageNum
    node := GetPage(cursor.Table.Pager, pn)

    cursor.CellNum = cursor.CellNum + 1
    if cursor.CellNum >= leafNodeNumCells(node) {
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
    PageNum uint32
    CellNum uint32
    EndOfTable bool
}

func StartTable(table *Table) *Cursor {

    c := &Cursor{
        Table: table,
        PageNum: table.RootPageNum,
        CellNum: 0,
    }

    // Root node
    rn := GetPage(table.Pager, table.RootPageNum)

    // # of cells
    nc := leafNodeNumCells(rn)

    if nc == 0 {
        c.EndOfTable = true
    } else {
        c.EndOfTable = false
    }

    return c

}

func EndTable(table *Table) *Cursor {

    c :=  &Cursor{
        Table: table,
        PageNum: table.RootPageNum,
    }

    // Root node
    rn := GetPage(table.Pager, table.RootPageNum)

    // # of cells
    nc := leafNodeNumCells(rn)

    c.CellNum = nc
    c.EndOfTable = true

    return c

}



