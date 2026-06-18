package main

import "fmt"

type StatementType int

const (
	StatementInsert StatementType = iota
	StatementSelect
)

type Statement struct {
	Type        StatementType
	RowToInsert Row
}

func PrepareStatement(input string, stmt *Statement) error {

	if input == "select" {
		stmt.Type = StatementSelect
		return nil
	}

	if len(input) >= 6 && input[:6] == "insert" {

		stmt.Type = StatementInsert

		_, err := fmt.Sscanf(
			input,
			"insert %d %s %s",
			&stmt.RowToInsert.ID,
			&stmt.RowToInsert.Username,
			&stmt.RowToInsert.Email,
		)

		return err
	}

	return fmt.Errorf("unrecognized statement")
}
