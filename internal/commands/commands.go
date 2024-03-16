package commands

import (
	"fmt"
	"os"

  "github.com/ammon134/pokedexcli/internal/locations"
)

type CliCommand struct {
	Callback    func() error
	Name        string
	Description string
}

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Display this help message.",
			Callback:    cmd_help,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex CLI.",
			Callback:    cmd_exit,
		},
		"map": {
			Name:        "map",
			Description: "Display the next 20 locations.",
			Callback:    cmd_exit,
		},
		"bmap": {
			Name:        "bmap",
			Description: "Display the previous 20 locations.",
			Callback:    cmd_exit,
		},
	}
}

func cmd_help() error {
	fmt.Println("Welcome to the Pokedex CLI!")
	fmt.Println("Usage: ")
	fmt.Println()
	for _, command := range GetCommands() {
		fmt.Printf("%s: %s\n", command.Name, command.Description)
	}
	return nil
}

func cmd_exit() error {
	os.Exit(0)
	return nil
}

func cmd_map() error {
	// Init the currentPage struct here
  cp := locations.CurrentPage{}
	return nil
}

func cmd_mapb() error {
	// Check if the currentPage struct is init or not,
	// if not throw the same error as is on the first page
	return nil
}
