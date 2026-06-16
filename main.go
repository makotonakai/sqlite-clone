package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printPrompt() {
	fmt.Print("db > ")
}

func main() {

	if len(os.Args) < 2 {

		fmt.Println(
			"Must supply a database filename.",
		)

		os.Exit(1)
	}

	table, err :=
		dbOpen(os.Args[1])

	if err != nil {
		panic(err)
	}

	defer dbClose(table)

	reader :=
		bufio.NewReader(os.Stdin)

	for {

		printPrompt()

		input, err :=
			reader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}

		input =
			strings.TrimSpace(input)

		if len(input) == 0 {
			continue
		}

		if input[0] == '.' {

			switch input {

			case ".exit":
				return

			case ".btree":

				fmt.Println("Tree:")

				printLeafNode(
					getPage(
						table.pager,
						table.rootPageNum,
					),
				)

				continue

			case ".constants":

				fmt.Println("Constants:")

				printConstants()

				continue

			default:

				fmt.Printf(
					"Unrecognized command '%s'\n",
					input,
				)

				continue
			}
		}

		var statement Statement

		err =
			prepareStatement(
				input,
				&statement,
			)

		if err != nil {

			fmt.Println(err)

			continue
		}

		err =
			executeStatement(
				&statement,
				table,
			)

		if err != nil {

			switch err.Error() {

			case "duplicate key":
				fmt.Println(
					"Error: Duplicate key.",
				)

			case "table full":
				fmt.Println(
					"Error: Table full.",
				)

			default:
				fmt.Println(err)
			}

			continue
		}

		fmt.Println("Executed.")
	}
}
