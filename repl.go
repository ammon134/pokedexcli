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

// TODO: take pokedex out of cliState
type cliState struct {
	nextLocations *string
	prevLocations *string
	cache         pokecache.Cache
	pokedex       pokeapi.Pokedex
}

type cliCommand struct {
	callback    func(args []string, cs *cliState) error
	name        string
	description string
}

func startRepl(cs *cliState) {
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

		err := command.callback(args, cs)
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch the pokemon with provided name.",
			callback:    cmdCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View details of caught pokemons",
			callback:    cmdInspect,
		},
	}
}

func cmdHelp(args []string, cs *cliState) error {
	fmt.Println("Welcome to the Pokedex CLI!")
	fmt.Println("Usage: ")
	fmt.Println()
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func cmdExit(args []string, cs *cliState) error {
	os.Exit(0)
	return nil
}

func cmdMap(args []string, cs *cliState) error {
	locationRes, err := pokeapi.GetLocations(cs.nextLocations, cs.cache)
	if err != nil {
		return errors.New("error getting locations")
	}

	cs.nextLocations = &locationRes.Next
	cs.prevLocations = &locationRes.Previous

	fmt.Println("---")
	for _, location := range locationRes.Results {
		fmt.Printf("%s\n", location.Name)
	}
	fmt.Println("---")

	return nil
}

func cmdMapB(args []string, cs *cliState) error {
	// On first page, api will return empty string for prev
	if cs.prevLocations == nil || *cs.prevLocations == "" {
		return errors.New("no previous map locations")
	}
	locationRes, err := pokeapi.GetLocations(cs.prevLocations, cs.cache)
	if err != nil {
		return errors.New("error getting locations")
	}

	cs.nextLocations = &locationRes.Next
	cs.prevLocations = &locationRes.Previous

	fmt.Println("---")
	for _, location := range locationRes.Results {
		fmt.Printf("%s\n", location.Name)
	}
	fmt.Println("---")

	return nil
}

func cmdExplore(args []string, cs *cliState) error {
	if len(args) == 0 {
		return errors.New("missing location to explore")
	}

	fmt.Printf("Exploring %v...\n", args[0])
	lDRes, err := pokeapi.GetLocationDetail(args[0], cs.cache)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemonEncounter := range lDRes.PokemonEncounters {
		fmt.Printf("- %v\n", pokemonEncounter.Pokemon.Name)
	}
	fmt.Println("---")

	return nil
}

func cmdCatch(args []string, cs *cliState) error {
	if len(args) == 0 {
		return errors.New("please provide pokemon name")
	}
	pokemonName := args[0]
	pokemon, err := pokeapi.GetPokemonData(pokemonName, cs.cache)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	caught, err := pokeapi.CatchPokemon(pokemon, cs.pokedex)
	if err != nil {
		return err
	}

	if caught {
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	fmt.Println("---")

	return nil
}

func cmdInspect(args []string, cs *cliState) error {
	if len(args) == 0 {
		return errors.New("please provide pokemon name")
	}
	pokemonName := args[0]
	pokemon, ok := cs.pokedex.Get(pokemonName)
	if !ok {
		return errors.New("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Print("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Print("Types:\n")
	for _, ptype := range pokemon.Types {
		fmt.Printf("  - %s\n", ptype.Type.Name)
	}
	fmt.Println("---")

	return nil
}
