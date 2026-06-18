package main

import (
	"bytes"
	"encoding/binary"
)

const (
	ColumnUsernameSize = 32
	ColumnEmailSize    = 255

	IDSize       = 4
	UsernameSize = ColumnUsernameSize
	EmailSize    = ColumnEmailSize

	RowSize = IDSize + UsernameSize + EmailSize
)

type Row struct {
	ID       uint32
	Username string
	Email    string
}

func SerializeRow(row *Row, dst []byte) {
	// ID
	binary.LittleEndian.PutUint32(
		dst[0:IDSize],
		row.ID,
	)

	// Fixed-width fields
	usernameField := dst[IDSize : IDSize+UsernameSize]

	emailField := dst[IDSize+UsernameSize : IDSize+UsernameSize+EmailSize]

	// Clear existing contents
	clear(usernameField)
	clear(emailField)

	// Copy string contents
	copy(usernameField, row.Username)
	copy(emailField, row.Email)
}

func DeserializeRow(src []byte) Row {
	usernameField := src[IDSize : IDSize+UsernameSize]

	emailField := src[IDSize+UsernameSize : IDSize+UsernameSize+EmailSize]

	return Row{
		ID: binary.LittleEndian.Uint32(
			src[0:IDSize],
		),

		Username: string(
			bytes.TrimRight(
				usernameField,
				"\x00",
			),
		),

		Email: string(
			bytes.TrimRight(
				emailField,
				"\x00",
			),
		),
	}
}
