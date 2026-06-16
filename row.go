package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	COLUMN_USERNAME_SIZE = 32
	COLUMN_EMAIL_SIZE    = 255

	ID_SIZE       = 4
	USERNAME_SIZE = COLUMN_USERNAME_SIZE
	EMAIL_SIZE    = COLUMN_EMAIL_SIZE

	ID_OFFSET       = 0
	USERNAME_OFFSET = ID_OFFSET + ID_SIZE
	EMAIL_OFFSET    = USERNAME_OFFSET + USERNAME_SIZE

	ROW_SIZE = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

type Row struct {
	ID       uint32
	Username string
	Email    string
}

func printRow(row *Row) {
	fmt.Printf("(%d, %s, %s)\n",
		row.ID,
		row.Username,
		row.Email)
}

func serializeRow(source *Row, destination []byte) {
	binary.LittleEndian.PutUint32(
		destination[ID_OFFSET:ID_OFFSET+ID_SIZE],
		source.ID,
	)

	copy(
		destination[
			USERNAME_OFFSET:
				USERNAME_OFFSET+USERNAME_SIZE],
		[]byte(source.Username),
	)

	copy(
		destination[
			EMAIL_OFFSET:
				EMAIL_OFFSET+EMAIL_SIZE],
		[]byte(source.Email),
	)
}

func deserializeRow(
	source []byte,
	destination *Row,
) {
	destination.ID =
		binary.LittleEndian.Uint32(
			source[ID_OFFSET : ID_OFFSET+ID_SIZE],
		)

	destination.Username =
		strings.TrimRight(
			string(
				source[
					USERNAME_OFFSET:
						USERNAME_OFFSET+USERNAME_SIZE],
			),
			"\x00",
		)

	destination.Email =
		strings.TrimRight(
			string(
				source[
					EMAIL_OFFSET:
						EMAIL_OFFSET+EMAIL_SIZE],
			),
			"\x00",
		)
}
