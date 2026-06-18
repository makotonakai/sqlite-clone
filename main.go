package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Must supply database filename.")
		os.Exit(1)
	}

	table, err := OpenTable(os.Args[1])

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {

		fmt.Print("db > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		if input == "" {
			continue
		}

		switch input {

		case ".exit":
				_ = table.Close()
				return

		case ".btree":
				root := table.Pager.GetPage(
						table.RootPageNum,
				)

				PrintLeafNode(root)
				continue

		case ".constants":
				PrintConstants()
				continue
		}

		var stmt Statement

		err := PrepareStatement(
			input,
			&stmt,
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		switch stmt.Type {

		case StatementInsert:

    row := stmt.RowToInsert

    cursor := TableFind(
      table,
      row.ID,
    )

    root := table.Pager.GetPage(
      table.RootPageNum,
    )

    numCells := LeafNodeNumCells(root)

    if cursor.CellNum < numCells {

        keyAtIndex := LeafNodeKey(
          root,
          cursor.CellNum,
        )

        if keyAtIndex == row.ID {
          fmt.Println(
            "Error: Duplicate key.",
          )
          continue
        }
    }

    err := LeafNodeInsert(
      cursor,
      row.ID,
      &row,
    )

    if err != nil {
        fmt.Println(err)
        continue
    }

    fmt.Println("Executed.")

		case StatementSelect:

			cursor := TableStart(table)

			for !cursor.EndOfTable {

				row :=
					DeserializeRow(
						cursor.Value(),
					)

				fmt.Printf(
					"(%d, %s, %s)\n",
					row.ID,
					row.Username,
					row.Email,
				)

				cursor.Advance()
			}

			fmt.Println("Executed.")
		}
	}
}
