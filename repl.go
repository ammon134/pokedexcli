package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ammon134/pokedexcli/internal/pokeapi"
	"github.com/ammon134/pokedexcli/internal/pokecache"
)

type cmdArg struct {
	nextLocations *string
	prevLocations *string
	cache         pokecache.Cache
}

type cliCommand struct {
	callback    func(ca *cmdArg, args []string) error
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
		command_list := strings.Fields(strings.ToLower(input))
		if len(command_list) == 0 {
			continue
		}

		command, ok := cmd_map[command_list[0]]
		if !ok {
			fmt.Println("not a command")
			continue
		}
		args := command_list[1:]

		err := command.callback(ca, args)
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
		"explore": {
			name:        "explore",
			description: "Show the pokemon in the provided location.",
			callback:    cmdExplore,
		},
	}
}

func cmdHelp(ca *cmdArg, args []string) error {
	fmt.Println("Welcome to the Pokedex CLI!")
	fmt.Println("Usage: ")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func cmdExit(ca *cmdArg, args []string) error {
	os.Exit(0)
	return nil
}

func cmdMap(ca *cmdArg, args []string) error {
	locationRes, err := pokeapi.GetLocations(ca.nextLocations, ca.cache)
	if err != nil {
		return errors.New("error getting locations")
	}

	ca.nextLocations = &locationRes.Next
	ca.prevLocations = &locationRes.Previous

	fmt.Println("---")
	for _, location := range locationRes.Results {
		fmt.Printf("%s\n", location.Name)
	}
	fmt.Println("---")

	return nil
}

func cmdMapB(ca *cmdArg, args []string) error {
	// On first page, api will return empty string for prev
	if ca.prevLocations == nil || *ca.prevLocations == "" {
		return errors.New("no previous map locations")
	}
	locationRes, err := pokeapi.GetLocations(ca.prevLocations, ca.cache)
	if err != nil {
		return errors.New("error getting locations")
	}

	ca.nextLocations = &locationRes.Next
	ca.prevLocations = &locationRes.Previous

	fmt.Println("---")
	for _, location := range locationRes.Results {
		fmt.Printf("%s\n", location.Name)
	}
	fmt.Println("---")

	return nil
}

func cmdExplore(ca *cmdArg, args []string) error {
	if len(args) == 0 {
		return errors.New("missing location to explore")
	}

	fmt.Printf("Exploring %v...\n", args[0])
	lDRes, err := pokeapi.GetLocationDetail(args[0], ca.cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemonEncounter := range lDRes.PokemonEncounters {
		fmt.Printf("- %v\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}
