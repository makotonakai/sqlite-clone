package main

type Cursor struct {
	Table      *Table
	PageNum    uint32
	CellNum    uint32
	EndOfTable bool
}

func TableStart(table *Table) *Cursor {

	cursor := &Cursor{
		Table:   table,
		PageNum: table.RootPageNum,
		CellNum: 0,
	}

	root := table.Pager.GetPage(0)

	cursor.EndOfTable =
		LeafNodeNumCells(root) == 0

	return cursor
}

func (c *Cursor) Value() []byte {

	page := c.Table.Pager.GetPage(c.PageNum)

	return LeafNodeValue(page, c.CellNum)
}

func (c *Cursor) Advance() {

	page := c.Table.Pager.GetPage(c.PageNum)

	c.CellNum++

	if c.CellNum >= LeafNodeNumCells(page) {
		c.EndOfTable = true
	}
}
