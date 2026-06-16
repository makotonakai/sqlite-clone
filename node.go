package main

import (
	"encoding/binary"
)

type NodeType uint8

const (
	NODE_INTERNAL NodeType = iota
	NODE_LEAF
)

/*
 * Common Node Header Layout
 */

const (
	NODE_TYPE_SIZE = 1
	NODE_TYPE_OFFSET = 0

	IS_ROOT_SIZE = 1
	IS_ROOT_OFFSET = NODE_TYPE_OFFSET + NODE_TYPE_SIZE

	PARENT_POINTER_SIZE = 4
	PARENT_POINTER_OFFSET =
		IS_ROOT_OFFSET + IS_ROOT_SIZE

	COMMON_NODE_HEADER_SIZE =
		NODE_TYPE_SIZE +
			IS_ROOT_SIZE +
			PARENT_POINTER_SIZE
)

/*
 * Leaf Node Header Layout
 */

const (
	LEAF_NODE_NUM_CELLS_SIZE = 4

	LEAF_NODE_NUM_CELLS_OFFSET =
		COMMON_NODE_HEADER_SIZE

	LEAF_NODE_HEADER_SIZE =
		COMMON_NODE_HEADER_SIZE +
			LEAF_NODE_NUM_CELLS_SIZE
)

/*
 * Leaf Node Body Layout
 */

const (
	LEAF_NODE_KEY_SIZE = 4
	LEAF_NODE_KEY_OFFSET = 0

	LEAF_NODE_VALUE_SIZE = ROW_SIZE

	LEAF_NODE_VALUE_OFFSET =
		LEAF_NODE_KEY_OFFSET +
			LEAF_NODE_KEY_SIZE

	LEAF_NODE_CELL_SIZE =
		LEAF_NODE_KEY_SIZE +
			LEAF_NODE_VALUE_SIZE
)

const (
	PAGE_SIZE = 4096
)

const (
	LEAF_NODE_SPACE_FOR_CELLS =
		PAGE_SIZE -
			LEAF_NODE_HEADER_SIZE

	LEAF_NODE_MAX_CELLS =
		LEAF_NODE_SPACE_FOR_CELLS /
			LEAF_NODE_CELL_SIZE
)

func leafNodeNumCells(node []byte) uint32 {
	return binary.LittleEndian.Uint32(
		node[
			LEAF_NODE_NUM_CELLS_OFFSET :
				LEAF_NODE_NUM_CELLS_OFFSET+
					LEAF_NODE_NUM_CELLS_SIZE],
	)
}

func setLeafNodeNumCells(
	node []byte,
	value uint32,
) {
	binary.LittleEndian.PutUint32(
		node[
			LEAF_NODE_NUM_CELLS_OFFSET :
				LEAF_NODE_NUM_CELLS_OFFSET+
					LEAF_NODE_NUM_CELLS_SIZE],
		value,
	)
}

func leafNodeCell(
	node []byte,
	cellNum uint32,
) []byte {
	offset :=
		LEAF_NODE_HEADER_SIZE +
			cellNum*LEAF_NODE_CELL_SIZE

	return node[offset:]
}

func leafNodeKey(
	node []byte,
	cellNum uint32,
) uint32 {

	offset :=
		LEAF_NODE_HEADER_SIZE +
			cellNum*LEAF_NODE_CELL_SIZE

	return binary.LittleEndian.Uint32(
		node[offset : offset+4],
	)
}

func setLeafNodeKey(
	node []byte,
	cellNum uint32,
	key uint32,
) {
	offset :=
		LEAF_NODE_HEADER_SIZE +
			cellNum*LEAF_NODE_CELL_SIZE

	binary.LittleEndian.PutUint32(
		node[offset:offset+4],
		key,
	)
}

func leafNodeValue(
	node []byte,
	cellNum uint32,
) []byte {

	offset :=
		LEAF_NODE_HEADER_SIZE +
			cellNum*LEAF_NODE_CELL_SIZE +
			LEAF_NODE_KEY_SIZE

	return node[offset : offset+ROW_SIZE]
}

func initializeLeafNode(
	node []byte,
) {
	node[NODE_TYPE_OFFSET] =
		byte(NODE_LEAF)

	setLeafNodeNumCells(node, 0)
}
