package main

import (
	"fmt"
)

type Table struct {
	pager       *Pager
	rootPageNum uint32
}

func dbOpen(filename string) (*Table, error) {

	pager, err := pagerOpen(filename)
	if err != nil {
		return nil, err
	}

	table := &Table{
		pager:       pager,
		rootPageNum: 0,
	}

	if pager.numPages == 0 {
		rootNode := getPage(pager, 0)
		initializeLeafNode(rootNode)
	}

	return table, nil
}

func dbClose(table *Table) error {

	pager := table.pager

	for i := uint32(0); i < pager.numPages; i++ {

		if pager.pages[i] == nil {
			continue
		}

		if err := pagerFlush(pager, i); err != nil {
			return err
		}

		pager.pages[i] = nil
	}

	return pager.file.Close()
}

func printConstants() {

	fmt.Printf("ROW_SIZE: %d\n", ROW_SIZE)
	fmt.Printf("COMMON_NODE_HEADER_SIZE: %d\n",
		COMMON_NODE_HEADER_SIZE)
	fmt.Printf("LEAF_NODE_HEADER_SIZE: %d\n",
		LEAF_NODE_HEADER_SIZE)
	fmt.Printf("LEAF_NODE_CELL_SIZE: %d\n",
		LEAF_NODE_CELL_SIZE)
	fmt.Printf("LEAF_NODE_SPACE_FOR_CELLS: %d\n",
		LEAF_NODE_SPACE_FOR_CELLS)
	fmt.Printf("LEAF_NODE_MAX_CELLS: %d\n",
		LEAF_NODE_MAX_CELLS)
}

func printLeafNode(node []byte) {

	numCells := leafNodeNumCells(node)

	fmt.Printf("leaf (size %d)\n", numCells)

	for i := uint32(0); i < numCells; i++ {
		fmt.Printf(
			"  - %d : %d\n",
			i,
			leafNodeKey(node, i),
		)
	}
}
