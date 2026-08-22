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

    LEAF_NODE_RIGHT_SPLIT_COUNT = (LEAF_NODE_MAX_CELLS + 1) / 2
    LEAF_NODE_LEFT_SPLIT_COUNT = (LEAF_NODE_MAX_CELLS + 1) - LEAF_NODE_RIGHT_SPLIT_COUNT

    INTERNAL_NODE_NUM_KEYS_SIZE = 4
    INTERNAL_NODE_NUM_KEYS_OFFSET = COMMON_NODE_HEADER_SIZE
    INTERNAL_NODE_RIGHT_CHILD_SIZE = 4
    INTERNAL_NODE_RIGHT_CHILD_OFFSET = INTERNAL_NODE_NUM_KEYS_OFFSET + INTERNAL_NODE_NUM_KEYS_SIZE
    INTERNAL_NODE_HEADER_SIZE = COMMON_NODE_HEADER_SIZE + INTERNAL_NODE_NUM_KEYS_SIZE + INTERNAL_NODE_RIGHT_CHILD_SIZE

    INTERNAL_NODE_KEY_SIZE = 4
    INTERNAL_NODE_CHILD_SIZE = 4
    INTERNAL_NODE_CELL_SIZE = INTERNAL_NODE_CHILD_SIZE + INTERNAL_NODE_KEY_SIZE

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
    nc := GetLeafNodeNumCells(node)
    fmt.Printf("leaf (size %d)\n", nc)

    for i := 0; i < int(nc); i++ {
        key := GetLeafNodeKey(node, uint32(i));
        fmt.Printf("  - %d : %d\n", uint32(i), key);
    }

}

func GetLeafNodeNumCells(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[LEAF_NODE_NUM_CELLS_OFFSET : LEAF_NODE_NUM_CELLS_OFFSET+4],
    )
}

func SetLeafNodeNumCells(node []byte, numCells uint32) {
    binary.LittleEndian.PutUint32(
        node[LEAF_NODE_NUM_CELLS_OFFSET : LEAF_NODE_NUM_CELLS_OFFSET+4],
        numCells,
    )
}

func GetLeafNodeCell(node []byte, cellNum uint32) []byte {
    offset := LEAF_NODE_HEADER_SIZE + int(cellNum)*LEAF_NODE_CELL_SIZE
    return node[offset : offset+LEAF_NODE_CELL_SIZE]
}

func GetLeafNodeKey(node []byte, cellNum uint32) uint32 {
    cell := GetLeafNodeCell(node, cellNum)

    return binary.LittleEndian.Uint32(
        cell[:LEAF_NODE_KEY_SIZE],
    )
}

func SetLeafNodeKey(node []byte, cellNum uint32, key uint32) {
    cell := GetLeafNodeCell(node, cellNum)

    binary.LittleEndian.PutUint32(
        cell[:LEAF_NODE_KEY_SIZE],
        key,
    )
}

func GetLeafNodeValue(node []byte, cellNum uint32) []byte {
    cell := GetLeafNodeCell(node, cellNum)

    return cell[LEAF_NODE_KEY_SIZE:]
}

func InitializeLeafNode(node []byte) {
    SetNodeType(node, NODE_LEAF)
    SetNodeRoot(node, false)
    SetLeafNodeNumCells(node, 0)
}

func InitializeInternalNode(node []byte) {
    SetNodeType(node, NODE_INTERNAL)
    SetNodeRoot(node, false)
    SetInternalNodeNumKeys(node, 0)
}

func GetInternalNodeNumKeys(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[INTERNAL_NODE_NUM_KEYS_OFFSET : INTERNAL_NODE_NUM_KEYS_OFFSET+4],
    )
}

func SetInternalNodeNumKeys(node []byte, key uint32) {
    binary.LittleEndian.PutUint32(
        node[INTERNAL_NODE_NUM_KEYS_OFFSET : INTERNAL_NODE_NUM_KEYS_OFFSET+4],
        key,
    )
}

func GetInternalNodeRightChild(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[INTERNAL_NODE_RIGHT_CHILD_OFFSET : INTERNAL_NODE_RIGHT_CHILD_OFFSET+INTERNAL_NODE_RIGHT_CHILD_SIZE],
    )
}

func SetInternalNodeRightChild(node []byte, child uint32) {
    binary.LittleEndian.PutUint32(
        node[INTERNAL_NODE_RIGHT_CHILD_OFFSET : INTERNAL_NODE_RIGHT_CHILD_OFFSET+INTERNAL_NODE_RIGHT_CHILD_SIZE],
        child,
    )
}

func GetInternalNodeCell(node []byte, cellNum uint32) []byte {
    offset := INTERNAL_NODE_HEADER_SIZE + int(cellNum)*INTERNAL_NODE_CELL_SIZE
    return node[offset : offset+INTERNAL_NODE_CELL_SIZE]
}

func GetInternalNodeChild(node []byte, cn uint32) uint32 {

    // Num keys
    nk := GetInternalNodeNumKeys(node)

    // Child num
    if cn > nk {
        fmt.Printf(
            "Tried to access child_num %d > num_keys %d\n",
            cn,
            nk,
        )
        os.Exit(1)
    }

    if cn == nk {
        return GetInternalNodeRightChild(node)
    }

    cell := GetInternalNodeCell(node, cn)

    return binary.LittleEndian.Uint32(
        cell[:INTERNAL_NODE_CHILD_SIZE],
    )
}

func SetInternalNodeChild(node []byte, childNum uint32, childPageNum uint32) {
    numKeys := GetInternalNodeNumKeys(node)

    if childNum > numKeys {
        fmt.Printf(
            "Tried to access child_num %d > num_keys %d\n",
            childNum,
            numKeys,
        )
        os.Exit(1)
    }

    if childNum == numKeys {
        SetInternalNodeRightChild(node, childPageNum)
        return
    }

    cell := GetInternalNodeCell(node, childNum)

    binary.LittleEndian.PutUint32(
        cell[:INTERNAL_NODE_CHILD_SIZE],
        childPageNum,
    )
}

func GetInternalNodeKey(node []byte, kn uint32) uint32 {
    cell := GetInternalNodeCell(node, kn)

    return binary.LittleEndian.Uint32(
        cell[INTERNAL_NODE_CHILD_SIZE:
            INTERNAL_NODE_CHILD_SIZE+INTERNAL_NODE_KEY_SIZE],
    )
}

func SetInternalNodeKey(node []byte, keyNum uint32, key uint32) {
    cell := GetInternalNodeCell(node, keyNum)

    binary.LittleEndian.PutUint32(
        cell[INTERNAL_NODE_CHILD_SIZE:
            INTERNAL_NODE_CHILD_SIZE+INTERNAL_NODE_KEY_SIZE],
        key,
    )
}

func GetNodeMaxKey(node []byte) uint32 {
    switch GetNodeType(node) {
    case NODE_INTERNAL:
        numKeys := GetInternalNodeNumKeys(node)
        return GetInternalNodeKey(node, numKeys-1)

    case NODE_LEAF:
        numCells := GetLeafNodeNumCells(node)
        return GetLeafNodeKey(node, numCells-1)
    }

    return 0
}

func InsertLeafNode(cursor *Cursor, key uint32, value *Row) {

    node := GetPage(cursor.Table.Pager, cursor.PageNum)

    // Number of cells currently in the node.
    nc := GetLeafNodeNumCells(node)

    // If the node is full, split it.
    if nc >= LEAF_NODE_MAX_CELLS {
        SplitAndInsertLeafNode(cursor, key, value)
        return
    }

    // Make room for the new cell by shifting cells to the right.
    if cursor.CellNum < nc {
        for i := int(nc); i > int(cursor.CellNum); i-- {
            copy(
                GetLeafNodeCell(node, uint32(i)),
                GetLeafNodeCell(node, uint32(i-1)),
            )
        }
    }

    // Increase number of cells.
    SetLeafNodeNumCells(node, nc+1)

    // Store key.
    SetLeafNodeKey(
        node,
        cursor.CellNum,
        key,
    )

    // Store row.
    SerializeRow(
        value,
        GetLeafNodeValue(node, cursor.CellNum),
    )
}

func FindLeafNode(table *Table, pageNum uint32, key uint32) *Cursor {

    node := GetPage(table.Pager, pageNum)

    // # of cells
    nc := GetLeafNodeNumCells(node)

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
        kai := GetLeafNodeKey(node, idx)

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


func GetNodeType(node []byte) NodeType {
    return NodeType(node[NODE_TYPE_OFFSET])
}

func SetNodeType(node []byte, nt NodeType) {
    node[NODE_TYPE_OFFSET] = byte(nt)
}

func SplitAndInsertLeafNode(cursor *Cursor, key uint32, value *Row) {

    // Old node
    on := GetPage(cursor.Table.Pager, cursor.PageNum)
    
    // New page number
    npn := GetUnusedPageNum(cursor.Table.Pager)

    // New node
    nn := GetPage(cursor.Table.Pager, npn)

    InitializeLeafNode(nn)

    for i := LEAF_NODE_MAX_CELLS; i >= 0; i-- {

        var dn []byte

        if i >= LEAF_NODE_LEFT_SPLIT_COUNT {
            // Destination node
            dn = nn
        } else {
            dn = on
        }

        // Index within node
        iwn := i % LEAF_NODE_LEFT_SPLIT_COUNT
        dst := GetLeafNodeCell(dn, uint32(iwn))

        if i == int(cursor.CellNum) {
            SerializeRow(value, dst)
        } else if i > int(cursor.CellNum) {
            copy(dst, GetLeafNodeCell(on, uint32(i - 1)))
        } else {
            copy(dst, GetLeafNodeCell(on, uint32(i)))
        }
    }

    SetLeafNodeNumCells(on, LEAF_NODE_LEFT_SPLIT_COUNT)
    SetLeafNodeNumCells(nn, LEAF_NODE_RIGHT_SPLIT_COUNT)



    if (IsNodeRoot(on)) {
        CreateNewRoot(cursor.Table, npn)
    } else {
        fmt.Printf("Need to implement updating parent after split\n")
        os.Exit(1)
    }

}

func CreateNewRoot(table *Table, rcpn uint32) {
    /*
    Handle splitting the root.
    Old root copied to new page, becomes left child.
    Address of right child passed in.
    Re-initialize root page to contain the new root node.
    New root node points to two children.
    */

    root := GetPage(table.Pager, table.RootPageNum)

    // right child
    // rc := GetPage(table.Pager, rcpn)

    // left child page num
    lcpn := GetUnusedPageNum(table.Pager)

    // left child
    lc := GetPage(table.Pager, lcpn)

    /* Left child has data copied from old root */
    // memcpy(left_child, root, PAGE_SIZE);
    copy(lc, root)

    // set_node_root(left_child, false);
    SetNodeRoot(lc, false)

    /* Root node is a new internal node with one key and two children */
    // initialize_internal_node(root);
    InitializeInternalNode(root)

    // set_node_root(root, true);
    SetNodeRoot(root, true)

    // *internal_node_num_keys(root) = 1;
    SetInternalNodeNumKeys(root, 1)

    // *internal_node_child(root, 0) = left_child_page_num;
    SetInternalNodeChild(root, 0, lcpn)

    // uint32_t left_child_max_key = get_node_max_key(left_child);
    // Left child max key 
    lcmk := GetNodeMaxKey(lc)

    // *internal_node_key(root, 0) = left_child_max_key;
    SetInternalNodeKey(root, 0, lcmk)

    // *internal_node_right_child(root) = right_child_page_num;
    SetInternalNodeRightChild(root, rcpn)
}

func IsNodeRoot(node []byte) bool {
    return node[IS_ROOT_OFFSET] != 0
}

func SetNodeRoot(node []byte, isRoot bool) {
    if isRoot {
        node[IS_ROOT_OFFSET] = 1
    } else {
        node[IS_ROOT_OFFSET] = 0
    }
}




