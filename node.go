package main

import (
	"fmt"
	"encoding/binary"
)

type NodeType uint8

const (
	NodeInternal NodeType = iota
	NodeLeaf
)

const (
	NodeTypeOffset = 0
	NodeTypeSize   = 1

	IsRootOffset = 1
	IsRootSize   = 1

	ParentOffset = 2
	ParentSize   = 4

	CommonNodeHeaderSize = 6
)

const (
	LeafNodeNumCellsOffset = CommonNodeHeaderSize
	LeafNodeNumCellsSize   = 4

	LeafNodeHeaderSize = CommonNodeHeaderSize + LeafNodeNumCellsSize
)

const (
	LeafNodeKeySize   = 4
	LeafNodeValueSize = RowSize

	LeafNodeCellSize = LeafNodeKeySize + LeafNodeValueSize

	LeafNodeSpaceForCells = PageSize - LeafNodeHeaderSize

	LeafNodeMaxCells = LeafNodeSpaceForCells / LeafNodeCellSize
)

const (
    InternalNodeNumKeysOffset = CommonNodeHeaderSize
    InternalNodeNumKeysSize   = 4

    InternalNodeRightChildOffset = InternalNodeNumKeysOffset + InternalNodeNumKeysSize
    InternalNodeRightChildSize = 4
    InternalNodeHeaderSize = CommonNodeHeaderSize + InternalNodeNumKeysSize + InternalNodeRightChildSize
)

const (
    InternalNodeChildSize = 4
    InternalNodeKeySize   = 4
    InternalNodeCellSize = InternalNodeChildSize + InternalNodeKeySize
)

const (
    LeafNodeRightSplitCount = (LeafNodeMaxCells + 1) / 2
    LeafNodeLeftSplitCount = (LeafNodeMaxCells + 1) - LeafNodeRightSplitCount
)

func GetNodeType(node []byte) NodeType {
	return NodeType(node[NodeTypeOffset])
}

func SetNodeType(node []byte, t NodeType) {
	node[NodeTypeOffset] = byte(t)
}

func InitializeLeafNode(node []byte) {
    SetNodeType(node, NodeLeaf)
    SetNodeRoot(node, false)
    SetLeafNodeNumCells(node, 0)
}

func LeafNodeNumCells(node []byte) uint32 {
	return binary.LittleEndian.Uint32(
		node[LeafNodeNumCellsOffset:],
	)
}

func SetLeafNodeNumCells(
	node []byte,
	value uint32,
) {
	binary.LittleEndian.PutUint32(
		node[LeafNodeNumCellsOffset:],
		value,
	)
}

func LeafNodeCell(
	node []byte,
	cellNum uint32,
) []byte {

	offset := LeafNodeHeaderSize + int(cellNum)*LeafNodeCellSize

	return node[offset:]
}

func LeafNodeKey(
	node []byte,
	cellNum uint32,
) uint32 {

	offset := LeafNodeHeaderSize + int(cellNum)*LeafNodeCellSize

	return binary.LittleEndian.Uint32(
		node[offset:],
	)
}

func SetLeafNodeKey(
	node []byte,
	cellNum uint32,
	key uint32,
) {

	offset := LeafNodeHeaderSize + int(cellNum)*LeafNodeCellSize

	binary.LittleEndian.PutUint32(
		node[offset:],
		key,
	)
}

func LeafNodeValue(
	node []byte,
	cellNum uint32,
) []byte {

	offset := LeafNodeHeaderSize + int(cellNum)*LeafNodeCellSize + LeafNodeKeySize

	return node[offset:]
}

func PrintLeafNode(node []byte) {
    numCells := LeafNodeNumCells(node)

    fmt.Printf(
        "leaf (size %d)\n",
        numCells,
    )

    for i := uint32(0); i < numCells; i++ {
        key := LeafNodeKey(node, i)

        fmt.Printf(
            "  - %d : %d\n",
            i,
            key,
        )
    }
}

func IsNodeRoot(node []byte) bool {
    return node[IsRootOffset] != 0
}

func SetNodeRoot(node []byte, isRoot bool) {
    if isRoot {
        node[IsRootOffset] = 1
    } else {
        node[IsRootOffset] = 0
    }
}

func InternalNodeNumKeys(node []byte) uint32 {
    return binary.LittleEndian.Uint32(
        node[InternalNodeNumKeysOffset:],
    )
}

func SetInternalNodeNumKeys(
    node []byte,
    value uint32,
) {
    binary.LittleEndian.PutUint32(
        node[InternalNodeNumKeysOffset:],
        value,
    )
}

func InternalNodeRightChild(
    node []byte,
) uint32 {

    return binary.LittleEndian.Uint32(
        node[InternalNodeRightChildOffset:],
    )
}

func SetInternalNodeRightChild(
    node []byte,
    child uint32,
) {
    binary.LittleEndian.PutUint32(
        node[InternalNodeRightChildOffset:],
        child,
    )
}

func InternalNodeCell(
    node []byte,
    cellNum uint32,
) []byte {

    offset :=
        InternalNodeHeaderSize +
            int(cellNum)*InternalNodeCellSize

    return node[offset:]
}

func InternalNodeChild(
    node []byte,
    childNum uint32,
) uint32 {

    numKeys :=
        InternalNodeNumKeys(node)

    if childNum > numKeys {
        panic("child out of range")
    }

    if childNum == numKeys {
        return InternalNodeRightChild(node)
    }

    return binary.LittleEndian.Uint32(
        InternalNodeCell(node, childNum),
    )
}

func InternalNodeKey(
    node []byte,
    keyNum uint32,
) uint32 {

    cell :=
        InternalNodeCell(node, keyNum)

    return binary.LittleEndian.Uint32(
        cell[InternalNodeChildSize:],
    )
}

func InitializeInternalNode(
    node []byte,
) {
    SetNodeType(node, NodeInternal)

    SetNodeRoot(node, false)

    SetInternalNodeNumKeys(
        node,
        0,
    )
}

func GetUnusedPageNum(
    pager *Pager,
) uint32 {
    return pager.NumPages
}

func GetNodeMaxKey(
    node []byte,
) uint32 {

    switch GetNodeType(node) {

    case NodeInternal:

        return InternalNodeKey(
            node,
            InternalNodeNumKeys(node)-1,
        )

    case NodeLeaf:

        return LeafNodeKey(
            node,
            LeafNodeNumCells(node)-1,
        )
    }

    panic("unknown node type")
}

func CreateNewRoot(
    table *Table,
    rightChildPageNum uint32,
) {

    root :=
        table.Pager.GetPage(
            table.RootPageNum,
        )

    leftChildPageNum :=
        GetUnusedPageNum(
            table.Pager,
        )

    leftChild :=
        table.Pager.GetPage(
            leftChildPageNum,
        )

    copy(leftChild, root)

    SetNodeRoot(
        leftChild,
        false,
    )

    InitializeInternalNode(root)

    SetNodeRoot(root, true)

    SetInternalNodeNumKeys(root, 1)

    SetInternalNodeChild(
        root,
        0,
        leftChildPageNum,
    )

    SetInternalNodeKeyValue(
        root,
        0,
        GetNodeMaxKey(leftChild),
    )

    SetInternalNodeRightChild(
        root,
        rightChildPageNum,
    )
}

func SetInternalNodeChild(
    node []byte,
    childNum uint32,
    child uint32,
) {

    numKeys := InternalNodeNumKeys(node)

    if childNum > numKeys {
        panic("child out of range")
    }

    if childNum == numKeys {

        SetInternalNodeRightChild(
            node,
            child,
        )

        return
    }

    cell := InternalNodeCell(
        node,
        childNum,
    )

    binary.LittleEndian.PutUint32(
        cell,
        child,
    )
}

func SetInternalNodeKeyValue(
    node []byte,
    keyNum uint32,
    key uint32,
) {

    cell := InternalNodeCell(
        node,
        keyNum,
    )

    binary.LittleEndian.PutUint32(
        cell[InternalNodeChildSize:],
        key,
    )
}

// func SetInternalNodeKeyValue(
//     node []byte,
//     keyNum uint32,
//     key uint32,
// ) {

//     cell := InternalNodeCell(
//         node,
//         keyNum,
//     )

//     binary.LittleEndian.PutUint32(
//         cell[InternalNodeChildSize:],
//         key,
//     )
// }

func LeafNodeSplitAndInsert(
    cursor *Cursor,
    key uint32,
    value *Row,
) {

    oldNode := cursor.Table.Pager.GetPage(
            cursor.PageNum,
        )

    newPageNum := GetUnusedPageNum(
            cursor.Table.Pager,
        )

    newNode := cursor.Table.Pager.GetPage(
            newPageNum,
        )

    InitializeLeafNode(newNode)

    for i := int(LeafNodeMaxCells); i >= 0; i-- {

        var destinationNode []byte

        if uint32(i) >= LeafNodeLeftSplitCount {
            destinationNode = newNode
        } else {
            destinationNode = oldNode
        }

        indexWithinNode := uint32(i) % LeafNodeLeftSplitCount

        destination := LeafNodeCell(
                destinationNode,
                indexWithinNode,
            )

        switch {

        case uint32(i) == cursor.CellNum:

            binary.LittleEndian.PutUint32(
                destination,
                key,
            )

            SerializeRow(
                value,
                destination[LeafNodeKeySize:],
            )

        case uint32(i) > cursor.CellNum:

            src := LeafNodeCell(
                  oldNode,
                  uint32(i-1),
                )

            copy(
                destination[:LeafNodeCellSize],
                src[:LeafNodeCellSize],
            )

        default:

            src := LeafNodeCell(
                    oldNode,
                    uint32(i),
                )

            copy(
                destination[:LeafNodeCellSize],
                src[:LeafNodeCellSize],
            )
        }
    }

    SetLeafNodeNumCells(
        oldNode,
        LeafNodeLeftSplitCount,
    )

    SetLeafNodeNumCells(
        newNode,
        LeafNodeRightSplitCount,
    )

    if IsNodeRoot(oldNode) {

        CreateNewRoot(
            cursor.Table,
            newPageNum,
        )

        return
    }

    panic(
        "Need to implement updating parent after split",
    )
}

func PrintTree(
    pager *Pager,
    pageNum uint32,
    level uint32,
) {

    node := pager.GetPage(pageNum)

    switch GetNodeType(node) {

    case NodeLeaf:

        numCells := LeafNodeNumCells(node)

        indent(level)
        fmt.Printf(
            "- leaf (size %d)\n",
            numCells,
        )

        for i := uint32(0); i < numCells; i++ {

            indent(level + 1)

            fmt.Printf(
                "- %d\n",
                LeafNodeKey(node, i),
            )
        }

    case NodeInternal:

        numKeys := InternalNodeNumKeys(node)

        indent(level)
        fmt.Printf(
            "- internal (size %d)\n",
            numKeys,
        )

        for i := uint32(0); i < numKeys; i++ {

            child :=
                InternalNodeChild(
                    node,
                    i,
                )

            PrintTree(
                pager,
                child,
                level+1,
            )

            indent(level + 1)

            fmt.Printf(
                "- key %d\n",
                InternalNodeKey(node, i),
            )
        }

        rightChild :=
            InternalNodeRightChild(node)

        PrintTree(
            pager,
            rightChild,
            level+1,
        )
    }
}

func indent(level uint32) {
    for i := uint32(0); i < level; i++ {
        fmt.Print("  ")
    }
}



