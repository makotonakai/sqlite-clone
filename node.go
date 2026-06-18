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

func GetNodeType(node []byte) NodeType {
	return NodeType(node[NodeTypeOffset])
}

func SetNodeType(node []byte, t NodeType) {
	node[NodeTypeOffset] = byte(t)
}

func InitializeLeafNode(node []byte) {

	SetNodeType(node, NodeLeaf)

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
