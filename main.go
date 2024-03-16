package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ammon134/pokedexcli/internal/commands"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cmd_map := commands.GetCommands()
	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		command_list := strings.Fields(input)
		if len(command_list) == 0 {
			continue
		}

		command, ok := cmd_map[command_list[0]]
		if !ok {
			cmd_map["help"].Callback()
			continue
		}

		command.Callback()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
