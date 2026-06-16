package main

import (
	"errors"
	"strconv"
	"strings"
)

type StatementType int

const (
	STATEMENT_INSERT StatementType = iota
	STATEMENT_SELECT
)

type Statement struct {
	statementType StatementType
	rowToInsert   Row
}

func prepareStatement(
	input string,
	statement *Statement,
) error {

	if strings.HasPrefix(input, "insert") {

		parts := strings.Fields(input)

		if len(parts) != 4 {
			return errors.New("syntax error")
		}

		id64, err := strconv.ParseUint(
			parts[1],
			10,
			32,
		)

		if err != nil {
			return errors.New("syntax error")
		}

		statement.statementType =
			STATEMENT_INSERT

		statement.rowToInsert = Row{
			ID:       uint32(id64),
			Username: parts[2],
			Email:    parts[3],
		}

		return nil
	}

	if input == "select" {

		statement.statementType =
			STATEMENT_SELECT

		return nil
	}

	return errors.New(
		"unrecognized statement",
	)
}

func leafNodeInsert(
	cursor *Cursor,
	key uint32,
	value *Row,
) {

	node :=
		getPage(
			cursor.table.pager,
			cursor.pageNum,
		)

	numCells :=
		leafNodeNumCells(node)

	if numCells >= LEAF_NODE_MAX_CELLS {

		panic(
			"Need to implement splitting a leaf node",
		)
	}

	if cursor.cellNum < numCells {

		for i := numCells; i > cursor.cellNum; i-- {

			dst :=
				leafNodeCell(
					node,
					i,
				)[:LEAF_NODE_CELL_SIZE]

			src :=
				leafNodeCell(
					node,
					i-1,
				)[:LEAF_NODE_CELL_SIZE]

			copy(dst, src)
		}
	}

	setLeafNodeNumCells(
		node,
		numCells+1,
	)

	setLeafNodeKey(
		node,
		cursor.cellNum,
		key,
	)

	serializeRow(
		value,
		leafNodeValue(
			node,
			cursor.cellNum,
		),
	)
}

func executeInsert(
	statement *Statement,
	table *Table,
) error {

	rootNode :=
		getPage(
			table.pager,
			table.rootPageNum,
		)

	if leafNodeNumCells(rootNode) >=
		LEAF_NODE_MAX_CELLS {

		return errors.New(
			"table full",
		)
	}

	if tableFindByID(
		table,
		statement.rowToInsert.ID,
	) {

		return errors.New(
			"duplicate key",
		)
	}

	cursor := tableEnd(table)

	leafNodeInsert(
		cursor,
		statement.rowToInsert.ID,
		&statement.rowToInsert,
	)

	return nil
}

func executeSelect(
	table *Table,
) {

	cursor := tableStart(table)

	var row Row

	for !cursor.endOfTable {

		deserializeRow(
			cursorValue(cursor),
			&row,
		)

		printRow(&row)

		cursorAdvance(cursor)
	}
}

func executeStatement(
	statement *Statement,
	table *Table,
) error {

	switch statement.statementType {

	case STATEMENT_INSERT:
		return executeInsert(
			statement,
			table,
		)

	case STATEMENT_SELECT:
		executeSelect(table)
		return nil
	}

	return nil
}
