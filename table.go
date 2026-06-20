package main

type Table struct {
	Pager       *Pager
	RootPageNum uint32
}

func OpenTable(filename string) (*Table, error) {
	pager, err := OpenPager(filename)
	if err != nil {
		return nil, err
	}

	table := &Table{
		Pager:       pager,
		RootPageNum: 0,
	}

	if pager.NumPages == 0 {
    root := pager.GetPage(0)
    InitializeLeafNode(root)
    SetNodeRoot(root, true)
	}

	return table, nil
}

func (t *Table) Close() error {
	return t.Pager.Close()
}

func TableFind(
	table *Table,
	key uint32,
) *Cursor {

	rootPageNum := table.RootPageNum

	rootNode := table.Pager.GetPage(rootPageNum)

	if GetNodeType(rootNode) == NodeLeaf {

		return LeafNodeFind(
			table,
			rootPageNum,
			key,
		)
	}

	return InternalNodeFind(
    table,
    rootPageNum,
    key,
	)
}

func LeafNodeFind(
	table *Table,
	pageNum uint32,
	key uint32,
) *Cursor {

	node := table.Pager.GetPage(pageNum)

	numCells := LeafNodeNumCells(node)

	cursor := &Cursor{
		Table:   table,
		PageNum: pageNum,
	}

	minIndex := uint32(0)
	onePastMax := numCells

	for minIndex != onePastMax {

		index := (minIndex + onePastMax) / 2

		keyAtIndex := LeafNodeKey(node, index)

		if key == keyAtIndex {

			cursor.CellNum = index
			return cursor
		}

		if key < keyAtIndex {
			onePastMax = index
		} else {
			minIndex = index + 1
		}
	}

	cursor.CellNum = minIndex

	return cursor
}

func LeafNodeInsert(
	cursor *Cursor,
	key uint32,
	value *Row,
) error {

	node :=
		cursor.Table.Pager.GetPage(
			cursor.PageNum,
		)

	numCells :=
		LeafNodeNumCells(node)

	if numCells >= LeafNodeMaxCells {
			LeafNodeSplitAndInsert(
			cursor,
			key,
			value,
		)

		return nil
	}

	if cursor.CellNum < numCells {

		for i := numCells; i > cursor.CellNum; i-- {

			dst :=
				LeafNodeCell(node, i)

			src :=
				LeafNodeCell(node, i-1)

			copy(
				dst[:LeafNodeCellSize],
				src[:LeafNodeCellSize],
			)
		}
	}

	SetLeafNodeNumCells(
		node,
		numCells+1,
	)

	SetLeafNodeKey(
		node,
		cursor.CellNum,
		key,
	)

	SerializeRow(
		value,
		LeafNodeValue(
			node,
			cursor.CellNum,
		),
	)

	return nil
}

func InternalNodeFind(
    table *Table,
    pageNum uint32,
    key uint32,
) *Cursor {

    node := table.Pager.GetPage(pageNum)

    numKeys := InternalNodeNumKeys(node)

    minIndex := uint32(0)
    maxIndex := numKeys

    for minIndex != maxIndex {

        index := (minIndex + maxIndex) / 2

        keyToRight := InternalNodeKey(
            node,
            index,
        )

        if key < keyToRight {
            maxIndex = index
        } else {
            minIndex = index + 1
        }
    }

    childNum := InternalNodeChild(
        node,
        minIndex,
    )

    child := table.Pager.GetPage(childNum)

    switch GetNodeType(child) {

    case NodeLeaf:
        return LeafNodeFind(
            table,
            childNum,
            key,
        )

    case NodeInternal:
        return InternalNodeFind(
            table,
            childNum,
            key,
        )
    }

    panic("unknown node type")
}
