package main

import (
	"time"

	"github.com/ammon134/pokedexcli/internal/pokeapi"
	"github.com/ammon134/pokedexcli/internal/pokecache"
)

func main() {
	cs := &cliState{
		nextLocations: nil,
		prevLocations: nil,
		cache:         pokecache.NewCache(5 * time.Minute),
		pokedex:       pokeapi.PokedexInit(),
	}
	startRepl(cs)
}
