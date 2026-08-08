package main

import (
    "os"
    "fmt"
	"encoding/binary"
)

type NodeType int

const (
	NODE_INTERNAL = iota
	NODE_LEAF
)

const (

	NODE_TYPE_SIZE = 1
	NODE_TYPE_OFFSET = 0

	IS_ROOT_SIZE = 1
	IS_ROOT_OFFSET = NODE_TYPE_SIZE

	PARENT_POINTER_SIZE = 4
	PARENT_POINTER_OFFSET = IS_ROOT_OFFSET + IS_ROOT_SIZE

	COMMON_NODE_HEADER_SIZE = NODE_TYPE_SIZE + IS_ROOT_SIZE + PARENT_POINTER_SIZE

	LEAF_NODE_NUM_CELLS_SIZE = 4
	LEAF_NODE_NUM_CELLS_OFFSET = COMMON_NODE_HEADER_SIZE
	LEAF_NODE_HEADER_SIZE = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELLS_SIZE

    LEAF_NODE_KEY_SIZE = 4
    LEAF_NODE_KEY_OFFSET = 0
    LEAF_NODE_VALUE_SIZE = ROW_SIZE
    LEAF_NODE_VALUE_OFFSET = LEAF_NODE_KEY_OFFSET + LEAF_NODE_KEY_SIZE
    LEAF_NODE_CELL_SIZE = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SIZE
    LEAF_NODE_SPACE_FOR_CELLS = PAGE_SIZE - LEAF_NODE_HEADER_SIZE
    LEAF_NODE_MAX_CELLS = LEAF_NODE_SPACE_FOR_CELLS / LEAF_NODE_CELL_SIZE

)

func PrintConstants() {
    fmt.Printf("ROW_SIZE: %d\n", ROW_SIZE);
    fmt.Printf("COMMON_NODE_HEADER_SIZE: %d\n", COMMON_NODE_HEADER_SIZE);
    fmt.Printf("LEAF_NODE_HEADER_SIZE: %d\n", LEAF_NODE_HEADER_SIZE);
    fmt.Printf("LEAF_NODE_CELL_SIZE: %d\n", LEAF_NODE_CELL_SIZE);
    fmt.Printf("LEAF_NODE_SPACE_FOR_CELLS: %d\n", LEAF_NODE_SPACE_FOR_CELLS);
    fmt.Printf("LEAF_NODE_MAX_CELLS: %d\n", LEAF_NODE_MAX_CELLS);
}

func PrintLeafNode(node []byte) {

    // # of cells
    nc := leafNodeNumCells(node)
    fmt.Printf("leaf (size %d)\n", nc)

    for i := 0; i < int(nc); i++ {
        key := leafNodeKey(node, uint32(i));
        fmt.Printf("  - %d : %d\n", uint32(i), key);
    }

}

func leafNodeNumCells(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[LEAF_NODE_NUM_CELLS_OFFSET : LEAF_NODE_NUM_CELLS_OFFSET+4],
    )
}

func setLeafNodeNumCells(node []byte, numCells uint32) {
    binary.LittleEndian.PutUint32(
        node[LEAF_NODE_NUM_CELLS_OFFSET : LEAF_NODE_NUM_CELLS_OFFSET+4],
        numCells,
    )
}

func leafNodeCell(node []byte, cellNum uint32) []byte {
    offset := LEAF_NODE_HEADER_SIZE + int(cellNum)*LEAF_NODE_CELL_SIZE
    return node[offset : offset+LEAF_NODE_CELL_SIZE]
}

func leafNodeKey(node []byte, cellNum uint32) uint32 {
    cell := leafNodeCell(node, cellNum)

    return binary.LittleEndian.Uint32(
        cell[:LEAF_NODE_KEY_SIZE],
    )
}

func setLeafNodeKey(node []byte, cellNum uint32, key uint32) {
    cell := leafNodeCell(node, cellNum)

    binary.LittleEndian.PutUint32(
        cell[:LEAF_NODE_KEY_SIZE],
        key,
    )
}

func leafNodeValue(node []byte, cellNum uint32) []byte {
    cell := leafNodeCell(node, cellNum)

    return cell[LEAF_NODE_KEY_SIZE:]
}

func initializeLeafNode(node []byte) {
    setNodeType(node, NODE_LEAF);
    setLeafNodeNumCells(node, 0)
}

func leafNodeInsert(cursor *Cursor, key uint32, value *Row) {

    node := GetPage(cursor.Table.Pager, cursor.PageNum)
    
    // # of cells
    nc := leafNodeNumCells(node)

    if nc >= LEAF_NODE_MAX_CELLS {
        fmt.Printf("Need to implement splitting a leaf node.\n")
        os.Exit(1)
    }

    if cursor.CellNum < nc {
        for i := int(nc); i >= int(cursor.CellNum); i-- {
            _ = copy(leafNodeCell(node, uint32(i)), leafNodeCell(node, uint32(i - 1)))
        }
    }

    setLeafNodeNumCells(node, leafNodeNumCells(node) + 1)
    setLeafNodeKey(node, cursor.CellNum, key)
    SerializeRow(value, leafNodeValue(node, cursor.CellNum))

}

func FindLeafNode(table *Table, pageNum uint32, key uint32) *Cursor {

    node := GetPage(table.Pager, pageNum)

    // # of cells
    nc := leafNodeNumCells(node)

    c := &Cursor{
        Table: table,
        PageNum: pageNum,
    }

    // Min index
    mi := uint32(0)

    // One past max index
    opmi := nc

    for opmi != mi {

        idx := (mi + opmi) / 2

        // Key at index
        kai := leafNodeKey(node, idx)

        if key == kai {
            c.CellNum = idx
            return c
        }

        if key < kai {
            opmi = idx
        } else {
            mi = idx + 1
        }
    }

    c.CellNum = mi
    return c
}


func nodeType(node []byte) NodeType {
    return NodeType(node[NODE_TYPE_OFFSET])
}

func setNodeType(node []byte, nt NodeType) {
    node[NODE_TYPE_OFFSET] = byte(nt)
}




