package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ammon134/pokedexcli/internal/locations"
)

type cmdArg struct {
	nextLocations *string
	prevLocations *string
}

type cliCommand struct {
	callback    func(ca *cmdArg) error
	name        string
	description string
}

func startRepl(ca *cmdArg) {
	scanner := bufio.NewScanner(os.Stdin)
	cmd_map := getCommands()
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
			fmt.Println("not a command")
			continue
		}

		err := command.callback(ca)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display this help message.",
			callback:    cmdHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex CLI.",
			callback:    cmdExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 locations.",
			callback:    cmdMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 locations.",
			callback:    cmdMapB,
		},
	}
}

func cmdHelp(ca *cmdArg) error {
	fmt.Println("Welcome to the Pokedex CLI!")
	fmt.Println("Usage: ")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func cmdExit(ca *cmdArg) error {
	os.Exit(0)
	return nil
}

func cmdMap(ca *cmdArg) error {
	locationRes, err := locations.GetLocations(ca.nextLocations)
	if err != nil {
		return errors.New("error getting locations")
	}

	ca.nextLocations = &locationRes.Next
	ca.prevLocations = &locationRes.Previous

	// On the "last" page, calling next will return
	// null for next
	if ca.nextLocations == nil {
		return errors.New("no next map locations")
	}

	for _, location := range locationRes.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func cmdMapB(ca *cmdArg) error {
	// On first page, api will return empty string for prev
	if ca.prevLocations == nil || *ca.prevLocations == "" {
		return errors.New("no previous map locations")
	}
	locationRes, err := locations.GetLocations(ca.prevLocations)
	if err != nil {
		return errors.New("error getting locations")
	}

	ca.nextLocations = &locationRes.Next
	ca.prevLocations = &locationRes.Previous

	for _, location := range locationRes.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}
