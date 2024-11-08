package main

import (
	"bufio"
	"fmt"
	"github.com/Nukambe/pokedex/commands"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	var code int

	for {
		// Read
		fmt.Print("gokedex > ")
		scanner.Scan()
		text = scanner.Text()
		fmt.Println("---------------")
		// Eval
		code = commands.ExecuteCommand(text)
		fmt.Println()
		if code < 0 {
			break
		}
	}
}
