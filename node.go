// node.go

package main

import (
    "os"
    "fmt"
    "math"
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

    LEAF_NODE_KEY_SIZE = 4
    LEAF_NODE_KEY_OFFSET = 0

    LEAF_NODE_VALUE_SIZE = ROW_SIZE
    LEAF_NODE_VALUE_OFFSET = LEAF_NODE_KEY_OFFSET + LEAF_NODE_KEY_SIZE

    LEAF_NODE_CELL_SIZE = LEAF_NODE_KEY_SIZE + LEAF_NODE_VALUE_SIZE
    LEAF_NODE_SPACE_FOR_CELLS = PAGE_SIZE - LEAF_NODE_HEADER_SIZE
    LEAF_NODE_MAX_CELLS = LEAF_NODE_SPACE_FOR_CELLS / LEAF_NODE_CELL_SIZE

    LEAF_NODE_RIGHT_SPLIT_COUNT = (LEAF_NODE_MAX_CELLS + 1) / 2
    LEAF_NODE_LEFT_SPLIT_COUNT = (LEAF_NODE_MAX_CELLS + 1) - LEAF_NODE_RIGHT_SPLIT_COUNT

    LEAF_NODE_NEXT_LEAF_SIZE = 4
    LEAF_NODE_NEXT_LEAF_OFFSET = LEAF_NODE_NUM_CELLS_OFFSET + LEAF_NODE_NUM_CELLS_SIZE
    LEAF_NODE_HEADER_SIZE = COMMON_NODE_HEADER_SIZE + LEAF_NODE_NUM_CELLS_SIZE + LEAF_NODE_NEXT_LEAF_SIZE


    INTERNAL_NODE_NUM_KEYS_SIZE = 4
    INTERNAL_NODE_NUM_KEYS_OFFSET = COMMON_NODE_HEADER_SIZE
    INTERNAL_NODE_RIGHT_CHILD_SIZE = 4
    INTERNAL_NODE_RIGHT_CHILD_OFFSET = INTERNAL_NODE_NUM_KEYS_OFFSET + INTERNAL_NODE_NUM_KEYS_SIZE
    INTERNAL_NODE_HEADER_SIZE = COMMON_NODE_HEADER_SIZE + INTERNAL_NODE_NUM_KEYS_SIZE + INTERNAL_NODE_RIGHT_CHILD_SIZE

    INTERNAL_NODE_KEY_SIZE = 4
    INTERNAL_NODE_CHILD_SIZE = 4
    INTERNAL_NODE_CELL_SIZE = INTERNAL_NODE_CHILD_SIZE + INTERNAL_NODE_KEY_SIZE
    INTERNAL_NODE_MAX_CELLS = 3

    INVALID_PAGE_NUM = math.MaxUint32

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

func GetLeafNodeNextLeaf(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[LEAF_NODE_NEXT_LEAF_OFFSET:LEAF_NODE_NEXT_LEAF_OFFSET+LEAF_NODE_NEXT_LEAF_SIZE],
    )
}

func SetLeafNodeNextLeaf(node []byte, nextLeaf uint32) {
    binary.LittleEndian.PutUint32(
        node[LEAF_NODE_NEXT_LEAF_OFFSET:LEAF_NODE_NEXT_LEAF_OFFSET+LEAF_NODE_NEXT_LEAF_SIZE],
        nextLeaf,
    )
}


func InitializeLeafNode(node []byte) {
    SetNodeType(node, NODE_LEAF)
    SetNodeRoot(node, false)
    SetLeafNodeNumCells(node, 0)
    SetLeafNodeNextLeaf(node, 0)
}

func InitializeInternalNode(node []byte) {
    SetNodeType(node, NODE_INTERNAL)
    SetNodeRoot(node, false)
    SetInternalNodeNumKeys(node, 0)
    SetInternalNodeRightChild(node, INVALID_PAGE_NUM)
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

func SetInternalNodeRightChild(node []byte, c uint32) {
    binary.LittleEndian.PutUint32(
        node[INTERNAL_NODE_RIGHT_CHILD_OFFSET : INTERNAL_NODE_RIGHT_CHILD_OFFSET+INTERNAL_NODE_RIGHT_CHILD_SIZE],
        c,
    )
}

func GetInternalNodeCell(node []byte, cellNum uint32) []byte {
    offset := INTERNAL_NODE_HEADER_SIZE + int(cellNum)*INTERNAL_NODE_CELL_SIZE
    return node[offset : offset+INTERNAL_NODE_CELL_SIZE]
}

func GetInternalNodeChild(node []byte, cn uint32) uint32 {
    nk := GetInternalNodeNumKeys(node)

    if cn > nk {
        fmt.Printf(
            "Tried to access child_num %d > num_keys %d\n",
            cn,
            nk,
        )
        os.Exit(1)
    }

    if cn == nk {
        rc := GetInternalNodeRightChild(node)

        if rc == INVALID_PAGE_NUM {
            fmt.Printf(
                "Tried to access right c of node, but was invalid page\n",
            )
            os.Exit(1)
        }

        return rc
    }

    cell := GetInternalNodeCell(node, cn)

    return binary.LittleEndian.Uint32(
        cell[:INTERNAL_NODE_CHILD_SIZE],
    )
}

func SetInternalNodeChild(node []byte, cn uint32, cpn uint32) {
    nk := GetInternalNodeNumKeys(node)

    if cn > nk {
        fmt.Printf(
            "Tried to access child_num %d > num_keys %d\n",
            cn,
            nk,
        )
        os.Exit(1)
    }

    if cn == nk {
        SetInternalNodeRightChild(node, cpn)
        return
    }

    cell := GetInternalNodeCell(node, cn)

    binary.LittleEndian.PutUint32(
        cell[:INTERNAL_NODE_CHILD_SIZE],
        cpn,
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

func FindInternalNodeChild(node []byte, key uint32) uint32 {
    // node := GetPage(table.Pager, pn)
    nk := GetInternalNodeNumKeys(node)

    minIndex := uint32(0)
    maxIndex := nk

    for minIndex != maxIndex {
        index := (minIndex + maxIndex) / 2

        keyToRight := GetInternalNodeKey(node, index)

        if keyToRight >= key {
            maxIndex = index
        } else {
            minIndex = index + 1
        }
    }

    return minIndex
    
}

func FindInternalNode(table *Table, pn uint32, key uint32) *Cursor {
    node := GetPage(table.Pager, pn)

    // Child index
    ci := FindInternalNodeChild(node, key)

    // Child num
    cn := GetInternalNodeChild(node, ci)

    // Child
    c := GetPage(table.Pager, cn)

    switch GetNodeType(c) {
    case NODE_LEAF:
        return FindLeafNode(table, cn, key)

    case NODE_INTERNAL:
        return FindInternalNode(table, cn, key)
    }

    return nil


}

func GetNodeMaxKey(pager *Pager, node []byte) uint32 {
    switch GetNodeType(node) {
    case NODE_INTERNAL:
        rightChild := GetInternalNodeRightChild(node)
        rightNode := GetPage(pager, rightChild)
        return GetNodeMaxKey(pager, rightNode)

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

func GetNodeParent(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[PARENT_POINTER_OFFSET : PARENT_POINTER_OFFSET+PARENT_POINTER_SIZE],
    )
}

func SetNodeParent(node []byte, ppn uint32) {
    binary.LittleEndian.PutUint32(
        node[PARENT_POINTER_OFFSET : PARENT_POINTER_OFFSET+PARENT_POINTER_SIZE],
        ppn,
    )
}

func SplitAndInsertLeafNode(cursor *Cursor, key uint32, value *Row) {

    // Old node
    on := GetPage(cursor.Table.Pager, cursor.PageNum)

    // Old max
    om := GetNodeMaxKey(cursor.Table.Pager, on)
    
    // New page number
    npn := GetUnusedPageNum(cursor.Table.Pager)

    // New node
    nn := GetPage(cursor.Table.Pager, npn)

    InitializeLeafNode(nn)
    SetNodeParent(nn, GetNodeParent(on))

    SetLeafNodeNextLeaf(nn, GetLeafNodeNextLeaf(on))
    SetLeafNodeNextLeaf(on, npn)

    for i := int(LEAF_NODE_MAX_CELLS); i >= 0; i-- {

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
            SetLeafNodeKey(dn, uint32(iwn), key)
            SerializeRow(
                value,
                GetLeafNodeValue(dn, uint32(iwn)),
            )
            SetLeafNodeKey(dn, uint32(iwn), key)
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
        // fmt.Printf("Need to implement updating p after split\n")
        // os.Exit(1)

        // Parent page num
        ppn := GetNodeParent(on)

        // New max
        nm := GetNodeMaxKey(cursor.Table.Pager, on)

        // Parent
        p := GetPage(cursor.Table.Pager, ppn)

        UpdateInternalNodeKey(p, om, nm)
        InsertInternalNode(cursor.Table, ppn, npn)
         
    }

}

func SplitAndInsertInternalNode(table *Table, ppn uint32, cpn uint32) {
    // Old page num
    opn := ppn

    // Old Node
    on := GetPage(table.Pager, opn)

    // Old max
    om := GetNodeMaxKey(table.Pager, on)

    // Child
    c := GetPage(table.Pager, cpn)

    // Child max
    cm := GetNodeMaxKey(table.Pager, c)

    // New page num
    npn := GetUnusedPageNum(table.Pager)

    // Splitting root
    sr := IsNodeRoot(on)

    var p []byte
    var nn []byte

    if sr {
        // The old root will become the left c.
        opn = GetInternalNodeChild(on, 0)
        on = GetPage(table.Pager, opn)

        // Parent
        p = GetPage(table.Pager, table.RootPageNum)

        // New node
        nn = GetPage(table.Pager, npn)
        InitializeInternalNode(nn)
    } else {
        p = GetPage(table.Pager, GetNodeParent(on))
        nn = GetPage(table.Pager, npn)
        InitializeInternalNode(nn)
    }

    // Old num keys
    onk := GetInternalNodeNumKeys(on)

    // Move the old right c into the new node.
    // Current page num
    cpn = GetInternalNodeRightChild(on)
    current := GetPage(table.Pager, cpn)

    InsertInternalNode(table, npn, cpn)
    SetNodeParent(current, npn)

    SetInternalNodeRightChild(on, INVALID_PAGE_NUM)

    // Move the children above the middle into the new node.
    for i := INTERNAL_NODE_MAX_CELLS - 1; i > INTERNAL_NODE_MAX_CELLS/2; i-- {
        cpn = GetInternalNodeChild(on, uint32(i))
        current = GetPage(table.Pager, cpn)

        InsertInternalNode(table, npn, cpn)
        SetNodeParent(current, npn)

        onk--
    }

    // The c immediately before the middle becomes
    // the old node's right c.
    SetInternalNodeRightChild(
        on,
        GetInternalNodeChild(on, onk-1),
    )
    onk--

    SetInternalNodeNumKeys(on, onk)

    // Decide which half receives the c being inserted.
    // Max after split
    mas := GetNodeMaxKey(table.Pager, on)

    // Destination page num
    dpn := opn

    if cm >= mas {
        dpn = npn
    }

    InsertInternalNode(table, dpn, cpn)
    SetNodeParent(c, dpn)

    // Update the p's key for the old node.
    UpdateInternalNodeKey(
        p,
        om,
        GetNodeMaxKey(table.Pager, on),
    )

    if sr {
        // The root already points to the old root's replacement structure.
        // CreateNewRoot has established the root.
    } else {
        // Parent page num
        ppn := GetNodeParent(on)

        InsertInternalNode(
            table,
            ppn,
            npn,
        )

        SetNodeParent(nn, ppn)
    }
}

func CreateNewRoot(table *Table, rcpn uint32) {
    /*
    Handle splitting the root.
    Old root copied to new page, becomes left c.
    Address of right c passed in.
    Re-initialize root page to contain the new root node.
    New root node points to two children.
    */

    root := GetPage(table.Pager, table.RootPageNum)

    // right c
    rc := GetPage(table.Pager, rcpn)

    // left c page num
    lcpn := GetUnusedPageNum(table.Pager)

    // left c
    lc := GetPage(table.Pager, lcpn)

    if GetNodeType(root) == NODE_INTERNAL {
        InitializeInternalNode(rc)
        InitializeInternalNode(lc)
    }

    /* Left c has data copied from old root */
    // memcpy(left_child, root, PAGE_SIZE);
    copy(lc, root)

    if GetNodeType(lc) == NODE_INTERNAL {
       var c []byte
       for i := 0; uint32(i) < GetInternalNodeNumKeys(lc); i++ {
            c = GetPage(table.Pager, GetInternalNodeChild(lc, uint32(i)))
            SetNodeParent(c, lcpn)
       } 
       c = GetPage(table.Pager, GetInternalNodeRightChild(lc))
       SetNodeParent(c, lcpn)
    }

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
    // Left c max key 
    lcmk := GetNodeMaxKey(table.Pager, lc)

    // *internal_node_key(root, 0) = left_child_max_key;
    SetInternalNodeKey(root, 0, lcmk)

    // *internal_node_right_child(root) = right_child_page_num;
    SetInternalNodeRightChild(root, rcpn)

    SetNodeParent(lc, table.RootPageNum)
    SetNodeParent(rc, table.RootPageNum)
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

func UpdateInternalNodeKey(node []byte, ok uint32, nk uint32) {
    // Old Child Index
    oci := FindInternalNodeChild(node, ok)
    SetInternalNodeKey(node, oci, nk)
}

func InsertInternalNode(table *Table, ppn uint32, cpn uint32) {

    // Parent
    p := GetPage(table.Pager, ppn)

    // Child
    c := GetPage(table.Pager, cpn)

    // Child max key
    cmk := GetNodeMaxKey(table.Pager, c)
    idx := FindInternalNodeChild(p, cmk)

    // Original Num Keys
    onk := GetInternalNodeNumKeys(p)

    if onk >= INTERNAL_NODE_MAX_CELLS {
        SplitAndInsertInternalNode(table, ppn, cpn)
        return
    }

    // Right c page num
    rcpn := GetInternalNodeRightChild(p)

    if rcpn == INVALID_PAGE_NUM {
        SetInternalNodeRightChild(p, cpn)
        return
    }

    // Right c
    rc := GetPage(table.Pager, rcpn)
    SetInternalNodeNumKeys(p, onk+1)

    if cmk > GetNodeMaxKey(table.Pager, rc) {
        SetInternalNodeChild(p, onk, rcpn)
        SetInternalNodeKey(p, onk, GetNodeMaxKey(table.Pager, rc))
        SetInternalNodeRightChild(p, cpn)
    } else {
        for i := 0; uint32(i) > idx; i-- {
            dst := GetInternalNodeCell(p, uint32(i))
            src := GetInternalNodeCell(p, uint32(i - 1))
            copy(dst, src)
        }
        SetInternalNodeChild(p, idx, cpn)
        SetInternalNodeKey(p, idx, cmk)
    }

}




