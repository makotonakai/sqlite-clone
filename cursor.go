package main

type Cursor struct {
	table      *Table
	pageNum    uint32
	cellNum    uint32
	endOfTable bool
}

func tableStart(
	table *Table,
) *Cursor {

	cursor := &Cursor{
		table:   table,
		pageNum: table.rootPageNum,
		cellNum: 0,
	}

	rootNode :=
		getPage(
			table.pager,
			table.rootPageNum,
		)

	numCells :=
		leafNodeNumCells(rootNode)

	cursor.endOfTable =
		numCells == 0

	return cursor
}

func tableEnd(
	table *Table,
) *Cursor {

	rootNode :=
		getPage(
			table.pager,
			table.rootPageNum,
		)

	numCells :=
		leafNodeNumCells(rootNode)

	return &Cursor{
		table:      table,
		pageNum:    table.rootPageNum,
		cellNum:    numCells,
		endOfTable: true,
	}
}

func cursorValue(
	cursor *Cursor,
) []byte {

	page :=
		getPage(
			cursor.table.pager,
			cursor.pageNum,
		)

	return leafNodeValue(
		page,
		cursor.cellNum,
	)
}

func cursorAdvance(
	cursor *Cursor,
) {

	page :=
		getPage(
			cursor.table.pager,
			cursor.pageNum,
		)

	cursor.cellNum++

	if cursor.cellNum >=
		leafNodeNumCells(page) {

		cursor.endOfTable = true
	}
}

func tableFindByID(
	table *Table,
	id uint32,
) bool {

	cursor := tableStart(table)
	defer func() {}()

	var row Row

	for !cursor.endOfTable {

		deserializeRow(
			cursorValue(cursor),
			&row,
		)

		if row.ID == id {
			return true
		}

		cursorAdvance(cursor)
	}

	return false
}
